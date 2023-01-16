package test

import (
	"luago/state"
	"luago/utils"
)

func TestMetatable(data []byte, chunkName string) {
	utils.OpenDebug = false
	ls := state.New()
	//注册Go的print函数到Lua的全局环境表里，lua编译后，会用GETTABUP指令把脚本里对应的函数与全局环境表里注册的Go函数对应起来
	ls.Register("print", print)
	ls.Register("getmetatable", getMetatable)
	ls.Register("setmetatable", setMetatable)
	ls.Load(data, chunkName, "b")
	//执行主函数
	ls.Call(0, 0)
}
