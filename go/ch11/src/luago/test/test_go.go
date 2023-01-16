package test

import (
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
