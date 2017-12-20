package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v2"
)

func main() {

	err := do()
	if err != nil {
		panic(err)
	}
}

func do() error {
	data := []int{}

	s := bufio.NewScanner(os.Stdin)

	columns := make(map[string]bool)

	start := time.Now()

	var cnt int
	for s.Scan() {
		m := make(map[interface{}]interface{})
		// println(s.Text())

		if err := parse([]byte(s.Text()), m); err != nil {
			fmt.Println("Error at", cnt, err)
			continue
		}

		v := m["cod"]
		n, err := strconv.Atoi(v.(string))
		if err != nil {
			panic(err)
		}
		data = append(data, n)
		// for z := range m {
		// 	k := z.(string)
		// 	columns[k] = true
		// }
		cnt++
		if cnt&(1<<20-1) == 0 {
			fmt.Println("processed ", cnt)
		}
	}
	fmt.Println("Imported", cnt, "rows; in", time.Since(start))
	for k := range columns {
		fmt.Println(k)
	}
	if err := s.Err(); err != nil {
		return err
	}

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		return err
	}

	ioutil.WriteFile("cod.clm", buf.Bytes(), 0666)

	return nil
}

func parse(d []byte, m map[interface{}]interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			err = fmt.Errorf("ups")
		}
	}()
	return yaml.Unmarshal(d, &m)
}
