package jsonparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSimple(t *testing.T) {
	iter := NewIterator(nil, strings.NewReader(` { "a" : "asdf" }`))
	iter.ReadObject(func(iter *Iterator, field string) bool {
		if field != "a" {
			t.FailNow()
		}
		if string(iter.ReadString()) != "asdf" {
			t.FailNow()
		}
		return true
	})
}

func TestParse10Strings(t *testing.T) {
	iter := NewIterator(nil, strings.NewReader(` { "a" : "asdfa", "b" : "asdfb","c" : "asdfc",	"d" : "asdfd",
		"e" : "asdfe","f" : "asdff","g" : "asdfg","h" : "asdfh","i" : "asdfi","j" : "asdfj","k" : "asdfk" }`))

	cnt := byte(0)
	success := iter.ReadObject(func(iter *Iterator, field string) bool {
		if field[0] != 'a'+cnt {
			t.FailNow()
		}
		if string(iter.ReadString()) != "asdf"+string([]byte{'a' + cnt}) {
			t.FailNow()
		}
		cnt++
		return true
	})

	if !success {
		t.FailNow()
	}
}

func TestJSONTestSuite(t *testing.T) {
	f, err := ioutil.ReadDir("test_parsing")
	if err != nil {
		t.FailNow()
	}

	for _, file := range f {
		res := testCase(filepath.Join("test_parsing", file.Name()))
		if res {
			fmt.Println(file.Name(), res)
		}
	}
}

func testCase(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	iter := NewIterator(nil, f)
	success := iter.ReadObject(func(iter *Iterator, field string) bool {
		iter.ReadString()
		return true
	})
	return success
}
