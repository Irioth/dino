package query

import (
	"fmt"
	"io"
	"testing"
)

type cl struct {
	name   string
	values []interface{}
}

func (c *cl) Name() string                { return c.name }
func (c *cl) Value(index int) interface{} { return c.values[index] }

type bs struct {
	values []int
	index  int
}

func (b *bs) Next() (int, error) {
	b.index++
	if b.index >= len(b.values) {
		return 0, io.EOF
	}
	return b.values[b.index], nil
}

func TestWorks(t *testing.T) {
	results := materialize(&bs{
		[]int{0, 2}, -1},
		&cl{"a", []interface{}{1, 2, 3}}, &cl{"b", []interface{}{"d", "e", "f"}})
	for i := 0; i < results.Columns(); i++ {
		fmt.Println(results.ColumnName(i))
	}
	err := results.Next()
	for err == nil {
		for i := 0; i < results.Columns(); i++ {
			fmt.Println(results.Value(i))
		}
		err = results.Next()
	}
	if err != io.EOF {
		panic(err)
	}
}
