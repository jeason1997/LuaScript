package test

import (
	"fmt"
	. "luago/api"
	"luago/state"
)

func TestGo(data []byte, chunkName string) {
	ls := state.New()
	//注册Go的print函数到Lua的全局环境表里，lua编译后，会用GETTABUP指令把脚本里的print函数与全局环境表里注册的Go函数对应起来
	ls.Register("print", print)
	ls.Load(data, chunkName, "b")
	//执行主函数
	ls.Call(0, 0)
}

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
