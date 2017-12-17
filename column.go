package dino

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type ColumnFactory struct {
}

func (f ColumnFactory) NewColumnData() columnData {
	return &Vector{}
}

type Vector struct {
	Data []interface{}
}

func (v *Vector) AddValue(index int, value interface{}) {
	for len(v.Data) < index {
		v.Data = append(v.Data, nil)
	}
	v.Data = append(v.Data, value)
}

func (v *Vector) Save(fname string) {
	if err := os.MkdirAll(filepath.Dir(fname), os.ModePerm); err != nil {
		panic(err)
	}
	e := NewEncoder()
	e.Encode(v)
	d, _ := e.Result()
	if err := ioutil.WriteFile(fname, d, 0666); err != nil {
		panic(err)
	}
}

func (v *Vector) Load(fname string) {
	data, _ := ioutil.ReadFile(fname)
	d := NewDecoder(data)
	d.Decode(v)
	d.Result()
}
