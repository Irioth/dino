package main

import (
	"os"

	"github.com/json-iterator/go"
)

func main() {

	it := jsoniter.Parse(jsoniter.ConfigDefault, os.Stdin, 65536)
	cnt := 0
	for it.ReadObjectCB(func(it *jsoniter.Iterator, field string) bool {
		it.ReadString()
		// println(field, it.WhatIsNext())
		// it.Skip()
		return true
	}) {
		cnt++
		if cnt&0xffff == 0 {
			println(cnt)
		}
		// println("-------------")
	}
	println(it.Error)

}
