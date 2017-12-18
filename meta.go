package dino

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const MetaFormatVersion = 2

var MetaMagic = []byte("DINOMT")

func (db *DB) loadMeta() error {
	fmeta := filepath.Join(db.path, "meta")
	dmeta, err := ioutil.ReadFile(fmeta)
	if err != nil {
		return err
	}

	// magic
	if !bytes.Equal(dmeta[:len(MetaMagic)], MetaMagic) {
		return errors.New("failed to read db meta, magic don't match: " + fmeta)
	}
	if dmeta[len(MetaMagic)] != MetaFormatVersion {
		return errors.New("unsupported meta version '" + strconv.Itoa(int(dmeta[len(MetaMagic)])) + "': " + fmeta)
	}
	dec := gob.NewDecoder(bytes.NewBuffer(dmeta[len(MetaMagic)+1:]))

	if err := dec.Decode(&db.tables); err != nil {
		return err
	}

	for _, t := range db.tables {
		t.factory = ColumnFactory{}
		for _, c := range t.columns {
			c.data = t.factory.NewColumnData()
			// c.data.Load(c.path)
		}

	}
	return nil
}

func (db *DB) saveMeta() error {
	var buf bytes.Buffer
	buf.Write(MetaMagic)                 // magic
	buf.Write([]byte{MetaFormatVersion}) // version

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(db.tables); err != nil {
		return err
	}
	fmeta := filepath.Join(db.path, "meta")
	return ioutil.WriteFile(fmeta, buf.Bytes(), os.ModePerm)
}

func (t *Table) GobEncode() ([]byte, error) {
	e := NewEncoder()
	e.Encode(t.name)
	e.Encode(t.path)
	e.Encode(t.size)
	e.Encode(t.columns)
	return e.Result()
}

func (t *Table) GobDecode(b []byte) error {
	d := NewDecoder(b)
	d.Decode(&t.name)
	d.Decode(&t.path)
	d.Decode(&t.size)
	d.Decode(&t.columns)
	return d.Result()
}

func (c *Column) GobEncode() ([]byte, error) {
	e := NewEncoder()
	e.Encode(c.name)
	e.Encode(c.path)
	return e.Result()
}
func (c *Column) GobDecode(d []byte) error {
	e := NewDecoder(d)
	e.Decode(&c.name)
	e.Decode(&c.path)
	return e.Result()
}

func (db *DB) DumpMeta(w io.Writer) {
	var totalRows, totalColumns int
	for _, t := range db.tables {
		fmt.Fprintf(w, "= %s ===============================\n", t.name)
		for _, c := range t.columns {
			fmt.Fprintf(w, "%s ", c.name)
		}
		fmt.Fprintf(w, "\nRows in table %d\n", t.RowsCount())
		totalRows += t.RowsCount()
		totalColumns += len(t.columns)
	}
	fmt.Fprintf(w, "================================\n")
	fmt.Fprintf(w, "Total: tables %d, columns %d, rows %d\n", len(db.tables), totalColumns, totalRows)
}
