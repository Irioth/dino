package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Irioth/dino"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli"
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

	time.Sleep(time.Minute)

	return nil
}

func importData(table *dino.Table, r io.Reader) error {

	fmt.Println("import data into", table.Name())

	it := jsoniter.Parse(jsoniter.ConfigFastest, os.Stdin, 65536)
	cnt := 0
	defer func() { println(cnt) }()
	for it.ReadObjectCB(func(it *jsoniter.Iterator, field string) bool {
		c := table.ColumnNew(field)
		value := it.ReadString()
		c.AddValue(table.RowsCount(), value)
		return true
	}) {
		table.IncRows()
		cnt++
		if cnt&0xffff == 0 {
			fmt.Println("Imported", cnt, "rows")
		}
		// println("-------------")
	}
	fmt.Println("Total imported", cnt)
	fmt.Println("Total table size:", table.RowsCount())

	return it.Error
}
