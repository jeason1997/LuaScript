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
			test.TestUndump(load_lua_data())
		case "2":
			test.TestStack()
		case "3":
			test.TestArith()
		case "4":
			test.TestVM(load_lua_data(), os.Args[2])
		case "5":
			test.TestGo(load_lua_data(), os.Args[2])
		case "6":
			test.TestMetatable(load_lua_data(), os.Args[2])
		}
	}
}

func load_lua_data() []byte {
	if len(os.Args) > 2 {
		data, err := ioutil.ReadFile(os.Args[2])

		if err != nil {
			panic(err)
		}

		return data
	}
	return nil
}
