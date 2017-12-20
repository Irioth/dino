package query

import "io"

type bitset interface {
	Next() (int, error)
}

type column interface {
	Name() string
	Value(index int) interface{}
}

type condition interface{}

type ResultStream interface {
	Next() error

	Columns() int
	ColumnName(index int) string

	Value(index int) interface{}
}

type mockresults struct {
	names  []string
	values [][]interface{}
	index  int
}

func newMockResults(names []string, values [][]interface{}) *mockresults {
	return &mockresults{names, values, -1}
}

func (f *mockresults) Next() error {
	f.index++
	if f.index >= len(f.values) {
		return io.EOF
	}
	return nil
}
func (f *mockresults) Columns() int {
	return len(f.names)
}
func (f *mockresults) ColumnName(index int) string {
	return f.names[index]
}
func (f *mockresults) Value(index int) interface{} {
	return f.values[f.index][index]
}

type emptyresults struct {
}

func (f *emptyresults) Next() error {
	return io.EOF
}
func (f *emptyresults) Columns() int {
	return 0
}
func (f *emptyresults) ColumnName(index int) string {
	return ""
}
func (f *emptyresults) Value(index int) interface{} {
	return nil
}

type lazyResults struct {
	current int
	filter  bitset
	columns []column
	err     error
}

func newLazyResults(filter bitset, columns ...column) *lazyResults {
	return &lazyResults{
		current: -1,
		filter:  filter,
		columns: columns,
	}
}

func (r *lazyResults) Next() error {
	if r.err != nil {
		return r.err
	}
	v, err := r.filter.Next()
	if err != nil {
		r.current = -1
		r.err = err
		return err
	}
	r.current = v
	return nil
}

func (r *lazyResults) Columns() int {
	return len(r.columns)
}

func (r *lazyResults) ColumnName(index int) string {
	return r.columns[index].Name()
}

func (r *lazyResults) Value(index int) interface{} {
	return r.columns[index].Value(r.current)
}

func materialize(filter bitset, columns ...column) ResultStream {
	return newLazyResults(filter, columns...)
	// return newMockResults([]string{"a", "b", "c"}, [][]interface{}{
	// 	{1, 2, 3},
	// 	{"d", "e", "f"},
	// })
}

func filter(cond condition, columns ...column) bitset {
	return nil
}

func do(op operation, columns ...column) column {

}
