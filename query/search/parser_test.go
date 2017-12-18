package search

import (
	"reflect"
	"testing"
)

func TestParseTable(t *testing.T) {
	q, err := ParseQuery("table")
	if err != nil {
		t.FailNow()
	}

	expected := &Query{table: "table"}

	if !reflect.DeepEqual(q, expected) {
		t.FailNow()
	}

}

func TestParseWhere(t *testing.T) {
	q, err := ParseQuery("table | where a=200")
	if err != nil {
		t.FailNow()
	}

	expected := &Query{table: "table", ops: []Operation{&Where{"a", EQUAL, "200"}}}

	if !reflect.DeepEqual(q, expected) {
		t.FailNow()
	}

}

func TestParseManyWhere(t *testing.T) {
	q, err := ParseQuery("table | where a=200|where b=40")
	if err != nil {
		t.FailNow()
	}

	expected := &Query{table: "table", ops: []Operation{&Where{"a", EQUAL, "200"}, &Where{"b", EQUAL, "40"}}}

	if !reflect.DeepEqual(q, expected) {
		t.FailNow()
	}

}
