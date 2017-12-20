package dino

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SmartColumn struct {
	Len        int
	Last       []string
	Compressed [][]byte
}

func (v *SmartColumn) len() int { return v.Len }

func (v *SmartColumn) add(s string) {
	v.Last = append(v.Last, s)
	v.Len++
	// println(v.Len)
	if len(v.Last) == 65536 {
		v.Compressed = append(v.Compressed, v.compress())
		v.Last = nil
	}
}

func (v *SmartColumn) compress() []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	for _, s := range v.Last {
		w.Write([]byte(s))
		w.Write([]byte{0x00})
	}
	w.Flush()
	w.Close()
	return b.Bytes()
}

func (v *SmartColumn) Get(index int) interface{} {
	return ""
}

func (v *SmartColumn) AddValue(index int, value interface{}) {
	for v.len() < index {
		v.add("")
	}
	v.add(value.(string))
}

func (v *SmartColumn) Save(fname string) {
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

func (v *SmartColumn) Load(fname string) {
	data, _ := ioutil.ReadFile(fname)
	d := NewDecoder(data)
	d.Decode(v)
	d.Result()
}
