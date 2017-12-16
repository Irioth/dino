package dino

import (
	"bytes"
	"encoding/gob"
)

type Encoder struct {
	b   bytes.Buffer
	g   *gob.Encoder
	err error
}

func NewEncoder() *Encoder {
	e := &Encoder{}
	e.g = gob.NewEncoder(&e.b)
	return e
}

func (e *Encoder) Encode(v interface{}) {
	if e.err != nil {
		return
	}
	e.err = e.g.Encode(v)
}

func (e *Encoder) Result() ([]byte, error) {
	return e.b.Bytes(), e.err
}

type Decoder struct {
	g   *gob.Decoder
	err error
}

func NewDecoder(d []byte) *Decoder {
	return &Decoder{gob.NewDecoder(bytes.NewReader(d)), nil}
}

func (d *Decoder) Decode(v interface{}) {
	if d.err != nil {
		return
	}
	d.err = d.g.Decode(v)
}

func (d *Decoder) Result() error {
	return d.err
}
