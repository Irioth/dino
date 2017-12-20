package dino

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

type StorableVector interface {
	Load(r io.Reader) error
	Save(w io.Writer) error
}

type VectorsRegistry struct {
	types map[string]reflect.Type
}

func NewVectorsRegistry() *VectorsRegistry {
	return &VectorsRegistry{make(map[string]reflect.Type)}
}

func (r *VectorsRegistry) Register(v StorableVector) {
	t := reflect.TypeOf(v).Elem()
	r.types[t.PkgPath()+"."+t.Name()] = t
}

func (r *VectorsRegistry) Load(rd io.Reader) (StorableVector, error) {
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	ln := int(binary.LittleEndian.Uint32(data))
	name := data[4 : 4+ln]
	typ := r.types[string(name)]
	vec := reflect.New(typ).Interface().(StorableVector)
	if err := vec.Load(bytes.NewReader(data[4+ln:])); err != nil {
		return nil, err
	}
	return vec, nil
}

func (r *VectorsRegistry) Save(w io.Writer, vec StorableVector) error {
	t := reflect.TypeOf(vec)
	name := t.Elem().PkgPath() + "." + t.Elem().Name()

	var x [4]byte
	binary.LittleEndian.PutUint32(x[:], uint32(len(name)))
	w.Write(x[:])
	w.Write([]byte(name))
	return vec.Save(w)
}

func (r *VectorsRegistry) LoadFile(fname string) (StorableVector, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return r.Load(bytes.NewReader(data))
}

func (r *VectorsRegistry) SaveFile(fname string, vec StorableVector) error {
	if err := os.MkdirAll(filepath.Dir(fname), 0666); err != nil {
		return err
	}
	var b bytes.Buffer
	if err := r.Save(&b, vec); err != nil {
		return err
	}
	return ioutil.WriteFile(fname, b.Bytes(), 0666)
}
