package test

import (
	"fmt"
	. "luago/api"
)

// Lua调用Go的print功能
func print(ls LuaState) int {
	//调用Go函数时，会新建一个调用帧压入调用帧列表，并且把参数压入到新的调用帧里
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print(ls.ToString(i))
		} else {
			fmt.Print(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

/*
 *Lua标准库提供了getmetatable()和setmetatable()函数，可以查询或者修改表的元表。
 *我们在本书第三部分会完整实现这两个函数，为了便于测试，这一章先实现简化版
 */
func getMetatable(ls LuaState) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func setMetatable(ls LuaState) int {
	ls.SetMetatable(1)
	return 1
}
