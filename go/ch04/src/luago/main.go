package main

import (
	"io/ioutil"
	"luago/test"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "1":
			test_luac()
		case "2":
			test_stack()
		}
	}
}

func test_luac() {
	if len(os.Args) > 2 {
		data, err := ioutil.ReadFile(os.Args[2])

		if err != nil {
			panic(err)
		}

		test.TestUndump(data)
	}
}

func test_stack() {
	test.TestStack()
}
