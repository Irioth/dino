package search

import (
	"fmt"

	"github.com/Irioth/dino"
)

type Executor struct {
	db *dino.DB
}

func NewExecutor(db *dino.DB) *Executor {
	return &Executor{db}
}

// TODO result type
func (e *Executor) Run(q *Query) (interface{}, error) {
	// TODO error on non existed table
	t := e.db.Table(q.table)
	for i := 0; i < t.RowsCount(); i++ {
		for _, column := range t.Columns() {
			t.Column(column).Get(i)
			// fmt.Printf("%s:%s ", column, t.Column(column).Get(i))
		}
		// fmt.Println()
	}
	fmt.Println("Total results:", t.RowsCount())
	return nil, nil
}
