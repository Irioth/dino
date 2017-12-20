package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	w := bufio.NewWriter(os.Stdout)

	s := bufio.NewScanner(os.Stdin)
	cnt := 0
	for s.Scan() {
		m := make(map[string]interface{})

		if err := yaml.Unmarshal([]byte(s.Text()), &m); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
			// return err
		}
		data, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		w.Write(data)
		w.WriteRune('\n')

		cnt++
		if cnt&0xffff == 0 {
			fmt.Fprintln(os.Stderr, "Imported", cnt, "rows")
		}
	}
	fmt.Fprintln(os.Stderr, "Imported", cnt, "rows")
	if err := s.Err(); err != nil {
		panic(err)
	}
	w.Flush()
}
