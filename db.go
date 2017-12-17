package dino

import (
	"os"
	"path/filepath"
)

type DB struct {
	path   string
	tables map[string]*Table
}

func newDB(path string) *DB {
	return &DB{
		path:   path,
		tables: make(map[string]*Table),
	}
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
	err := db.saveMeta()
	if err != nil {
		return err
	}
	for _, t := range db.tables {
		for _, c := range t.columns {
			c.data.Save(c.path)
		}
	}
	return nil
}

// Table returns existed table in Database or create new one with defualt params
func (db *DB) Table(name string) *Table {
	if t, ok := db.tables[name]; ok {
		return t
	}
	t := newTable(name, filepath.Join(db.path, name), ColumnFactory{})
	db.tables[name] = t
	return t
}
