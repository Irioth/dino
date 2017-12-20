package dino

import (
	"bytes"
	"io"
	"testing"
)

type fakeStorableVector struct {
}

func (v *fakeStorableVector) Save(w io.Writer) error {
	return nil
}

func (v *fakeStorableVector) Load(r io.Reader) error {
	return nil
}

func TestSaveVector(t *testing.T) {
	r := NewVectorsRegistry()
	r.Register((*fakeStorableVector)(nil))
	var b bytes.Buffer
	if err := r.Save(&b, &fakeStorableVector{}); err != nil {
		t.FailNow()
	}

	// fmt.Println(hex.Dump(b.Bytes()))

	vec, err := r.Load(&b)
	if err != nil {
		t.FailNow()
	}
	if _, ok := vec.(*fakeStorableVector); !ok {
		t.FailNow()
	}

}
