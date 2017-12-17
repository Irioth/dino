package dino

import (
	"fmt"
	"io"
	"path/filepath"
)

type Column struct {
	name string
	path string

	data columnData
}

type columnData interface {
	AddValue(index int, value interface{})
	Load(fname string)
	Save(fname string)
}

func (c *Column) AddValue(index int, value interface{}) {
	c.data.AddValue(index, value)
}

func (c *Column) Dump(w io.Writer) {
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
