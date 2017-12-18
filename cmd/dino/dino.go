package main

import (
	"fmt"
	"os"

	"github.com/Irioth/dino"
	"github.com/Irioth/dino/query/search"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "dino"
	app.Usage = "dino database"
	app.Version = "0.1"

	app.Commands = []cli.Command{
		{
			Name:      "init",
			Usage:     "create empty database",
			ArgsUsage: "<path>",
			Action:    initdb,
		},

		{
			Name:      "info",
			Usage:     "dump database meta",
			ArgsUsage: "<path>",
			Action:    info,
		},

		{
			Name:      "import",
			Usage:     "import data to table",
			ArgsUsage: "<path> <table>",
			Action:    importAction,
		},
		{
			Name:      "search",
			Usage:     "search over database",
			ArgsUsage: "<path> <query>",
			Action:    searchAction,
		},
	}

	app.Run(os.Args)
}

func initdb(c *cli.Context) error {
	path := c.Args().First()
	if path == "" {
		path = "."
		return cli.NewExitError("path must be specified", 1)
	}
	db, err := dino.Create(path)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err = db.Close()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func info(c *cli.Context) (err error) {
	path := c.Args().First()
	if path == "" {
		return cli.NewExitError("path must be specified", 1)
	}
	db, err := dino.Open(path)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			return
		}
	}()

	db.DumpMeta(os.Stdout)

	return nil
}

func searchAction(c *cli.Context) (err error) {
	path := c.Args().First()
	if path == "" {
		return cli.NewExitError("path must be specified", 1)
	}
	db, err := dino.Open(path)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			return
		}
	}()

	query, err := search.ParseQuery(c.Args()[1])
	if err != nil {
		fmt.Println(err)
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("%# v", query)

	return nil
}
