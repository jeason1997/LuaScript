// 测试虚拟机运算操作
package test

import (
	. "luago/api"
	"luago/state"
	"luago/utils"
)

func TestArith() {
	ls := state.New()
	ls.PushInteger(1)
	ls.PushString("2.0")
	ls.PushString("3.0")
	ls.PushNumber(4.0)
	utils.PrintStack(ls)

	ls.Arith(LUA_OPADD)
	utils.PrintStack(ls)
	ls.Arith(LUA_OPBNOT)
	utils.PrintStack(ls)
	ls.Len(2)
	utils.PrintStack(ls)
	ls.Concat(3)
	utils.PrintStack(ls)
	ls.PushBoolean(ls.Compare(1, 2, LUA_OPEQ))
	utils.PrintStack(ls)
}
