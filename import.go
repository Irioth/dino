package dino

import (
	"bufio"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

func (db *DB) Import(table string, r io.Reader) error {
	fmt.Println("import data into", table)
	db.AddTable(table)
	s := bufio.NewScanner(r)
	var cnt int
	for s.Scan() {
		m := make(map[interface{}]interface{})
		// println(s.Text())

		if err := yaml.Unmarshal([]byte(s.Text()), &m); err != nil {
			return err
		}
		for z := range m {
			k := z.(string)
			if _, ok := db.tables[table].Columns[k]; !ok {
				db.tables[table].Columns[k] = column{}
				fmt.Println("added column", k)
			}
		}
		cnt++
	}
	fmt.Println("Imported", cnt, "rows")
	if err := s.Err(); err != nil {
		return err
	}
	return nil
}
