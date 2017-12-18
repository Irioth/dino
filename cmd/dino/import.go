package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Irioth/dino"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

func importAction(c *cli.Context) (err error) {
	path := c.Args().First()
	if path == "" {
		return cli.NewExitError("path must be specified", 1)
	}
	db, err := dino.Open(path)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer func() { err = db.Close() }()

	tName := c.Args().Get(1)
	table := db.Table(tName)

	if err := importData(table, os.Stdin); err != nil {
		fmt.Println("error", err.Error())
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func importData(table *dino.Table, r io.Reader) error {
	fmt.Println("import data into", table.Name())
	s := bufio.NewScanner(r)
	cnt := 0
	for s.Scan() {
		m := make(map[string]interface{})

		if err := yaml.Unmarshal([]byte(s.Text()), &m); err != nil {
			return err
		}
		table.AppendRow(m)
		cnt++
		if cnt&0xffff == 0 {
			fmt.Println("Imported", cnt, "rows")
		}
	}
	fmt.Println("Imported", cnt, "rows")
	fmt.Println("Total size", table.RowsCount())
	if err := s.Err(); err != nil {
		return err
	}
	// table.Column("api").Dump(os.Stdout)
	return nil
}
