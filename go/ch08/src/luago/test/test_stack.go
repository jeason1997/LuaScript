// 测试虚拟机栈操作
package test

import (
	"luago/state"
	debugger "luago/utils"
)

func TestStack() {
	ls := state.New()
	ls.PushBoolean(true)
	debugger.PrintStack(ls)
	ls.PushInteger(10)
	debugger.PrintStack(ls)
	ls.PushNil()
	debugger.PrintStack(ls)
	ls.PushString("hello")
	debugger.PrintStack(ls)
	ls.PushValue(-4)
	debugger.PrintStack(ls)
	ls.Replace(3)
	debugger.PrintStack(ls)
	ls.SetTop(6)
	debugger.PrintStack(ls)
	ls.Remove(-3)
	debugger.PrintStack(ls)
	ls.SetTop(-5)
	debugger.PrintStack(ls)
}
