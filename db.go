package dino

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type column struct {
	Fname string
}

type table struct {
	Columns map[string]column
}

func newTable() table {
	return table{make(map[string]column)}
}

type DB struct {
	path   string
	tables map[string]table
}

func newDB(path string) *DB {
	return &DB{path: path, tables: make(map[string]table)}
}

func (db *DB) Dump() string {
	return fmt.Sprintf("%#v", db)
}

// Create empty Database
func Create(path string) (*DB, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}
	db := newDB(path)
	if err := db.saveMeta(); err != nil {
		return nil, err
	}
	return db, nil
}

// Open existing Database
func Open(path string) (*DB, error) {
	db := newDB(path)
	if err := db.loadMeta(); err != nil {
		return nil, err
	}
	return db, nil
}

// Close opened database
func (db *DB) Close() error {
	return db.saveMeta()
}

func (db *DB) CreateTable(name string) {
	db.tables[name] = newTable()
}

const MetaFormatVersion = 1

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
