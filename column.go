package dino

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ColumnFactory struct {
}

func (f ColumnFactory) NewColumnData() columnData {
	// return &Vector{}
	return &SmartColumn{}
	// return &StringVector{Counts: make(map[string]int), Indexes: make(map[string]int)}
}

func (f *StringVector) Add(v string) {
	if !f.Many {
		f.Counts[v] = f.Counts[v] + 1
		if p, ok := f.Indexes[v]; ok {
			f.Opt = append(f.Opt, byte(p))
		} else {
			if len(f.Indexes) < 256 {
				p := len(f.OptData)
				f.Indexes[v] = p
				f.Opt = append(f.Opt, byte(p))
				f.OptData = append(f.OptData, v)
			} else {
				f.Many = true
				fmt.Println("deopt")
				// deopt

				for i := range f.Opt {
					f.Data = append(f.Data, f.OptData[f.Opt[i]])
				}

				fmt.Println("deopt done")
				f.Indexes = nil
				f.Counts = nil
				f.OptData = nil
				f.Opt = nil
			}
		}
	} else {
		f.Data = append(f.Data, v)
	}
}

func (f *StringVector) Len() int {
	if f.Many {
		return len(f.Data)
	}
	return len(f.Opt)
}

type StringVector struct {
	Many    bool
	Counts  map[string]int
	Indexes map[string]int
	OptData []string
	Opt     []byte
	Data    []string
}

func (v *StringVector) Get(index int) interface{} {
	if index >= v.Len() {
		return ""
	}
	if v.Many {
		return v.Data[index]
	}
	return v.OptData[v.Opt[index]]
}

func (v *StringVector) AddValue(index int, value interface{}) {
	return
	for v.Len() < index {
		v.Add("")
	}
	v.Add(value.(string))
}

func (v *StringVector) Save(fname string) {
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

func (v *StringVector) Load(fname string) {
	v.Counts = make(map[string]int)
	v.Indexes = make(map[string]int)
	data, _ := ioutil.ReadFile(fname)
	d := NewDecoder(data)
	d.Decode(v)
	d.Result()

	fmt.Printf("\n%t %d\n%# v", v.Many, len(v.Counts), v.Counts)
}

type Vector struct {
	Data []interface{}
}

func (v *Vector) Get(index int) interface{} {
	if index >= len(v.Data) {
		return nil
	}
	return v.Data[index]
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
