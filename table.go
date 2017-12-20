package dino

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"unsafe"
)

type Column struct {
	name            string
	path            string
	loaded, changed bool

	data columnData
}

type columnData interface {
	AddValue(index int, value interface{})
	Get(index int) interface{}
	Load(fname string)
	Save(fname string)
}

func (c *Column) checkLoaded() {
	if !c.loaded {
		c.data.Load(c.path)
		c.loaded = true
	}
}

func (c *Column) Save() {
	if c.changed {
		c.data.Save(c.path)
		c.changed = false
	}
}

func (c *Column) AddValue(index int, value interface{}) {
	c.checkLoaded()
	c.data.AddValue(index, value)
	c.changed = true
}

func (c *Column) Get(index int) interface{} {
	c.checkLoaded()
	return c.data.Get(index)
}

func (c *Column) Dump(w io.Writer) {
	c.checkLoaded()
	fmt.Fprintf(w, "%# v", c.data)
}

type columnDataFactory interface {
	NewColumnData() columnData
}

type Table struct {
	name    string
	path    string
	columns map[string]*Column
	size    int
	factory columnDataFactory
}

func newTable(name, path string, factory columnDataFactory) *Table {
	return &Table{
		name:    name,
		path:    path,
		columns: make(map[string]*Column),
		size:    0,
		factory: factory,
	}
}

func (t *Table) Name() string               { return t.name }
func (t *Table) RowsCount() int             { return t.size }
func (t *Table) Column(name string) *Column { return t.columns[name] }
func (t *Table) ColumnNew(name string) *Column {
	if c, ok := t.columns[name]; ok {
		return c
	}
	return t.сreateColumn(dup(name))
}

func dup(x string) string {
	data := unsafeStrToByte(x)
	return string(data)
}

func unsafeStrToByte(s string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	var b []byte
	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	byteHeader.Data = strHeader.Data

	// need to take the length of s here to ensure s is live until after we update b's Data
	// field since the garbage collector can collect a variable once it is no longer used
	// not when it goes out of scope, for more details see https://github.com/golang/go/issues/9046
	l := len(s)
	byteHeader.Len = l
	byteHeader.Cap = l
	return b
}

func (t *Table) Columns() []string {
	var s []string
	for c := range t.columns {
		s = append(s, c)
	}
	return s
}

func (t *Table) AppendRow(data map[string]interface{}) {
	for k, v := range data {
		c, ok := t.columns[k]
		if !ok {
			c = t.сreateColumn(k)
			fmt.Println("add new column", k)
		}
		c.AddValue(t.size, v)
	}
	t.size++
}

func (t *Table) IncRows() {
	t.size++
}

func (t *Table) сreateColumn(name string) *Column {
	if _, ok := t.columns[name]; ok {
		panic("column already exists")
	}
	c := &Column{
		name: name,
		path: filepath.Join(t.path, name),
		data: t.factory.NewColumnData(),
	}
	t.columns[name] = c
	return c
}
