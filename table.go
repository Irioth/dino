package dino

import (
	"bytes"
	"encoding/gob"
	"path/filepath"
)

type column struct {
	name string
	path string
}

func (c *column) AddValue(index int, value interface{}) {

}

type Table struct {
	name    string
	path    string
	columns map[string]*column
	size    int
}

func (t *Table) GobEncode() ([]byte, error) {
	var x struct {
		name    string
		path    string
		columns map[string]*column
		size    int
	} = *t
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(&x); err != nil {
		return err
	}
	return b.Bytes(), nil
}

func newTable(name, path string) *Table {
	return &Table{name, path, make(map[string]*column), 0}
}

func (t *Table) Name() string { return t.name }
func (t *Table) Size() int    { return t.size }

func (t *Table) AppendRow(data map[string]interface{}) {
	for k, v := range data {
		c, ok := t.columns[k]
		if !ok {
			c = t.сreateColumn(k)
		}
		c.AddValue(t.size, v)
	}
	t.size++
}

func (t *Table) сreateColumn(name string) *column {
	if _, ok := t.columns[name]; ok {
		panic("column already exists")
	}
	t.columns[name] = &column{name, filepath.Join(t.path, name)}
}
