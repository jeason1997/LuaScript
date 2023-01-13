package state

import (
	. "luago/api"
	"luago/binchunk"
)

/*
*所谓闭包，就是按词法作用域捕获了非局部变量的嵌套函数。
*现在大家知道为什么在Lua内部函数被称为闭包了吧？因为Lua函数本质上全都是闭包。
*就算是编译器为我们生成的主函数也不例外，它从外部捕获了_ENV变量

*我们使用closure结构体来统一表示Lua和Go闭包。
*如果proto字段不是nil，说明这是Lua闭包。否则，goFunc字段一定不是nil，说明这是Go闭包。
 */
type closure struct {
	proto  *binchunk.Prototype //Lua闭包
	goFunc GoFunction          //Go闭包
	upvals []*upvalue          //闭包内部捕获的非局部变量（来自外围函数）
}

// Upvalue就是闭包内部捕获的非局部变量（来自外围函数）
// Lua编译器在生成主函数时会在它的外围隐式声明一个局部变量'_ENV'（存在注册表的全局环境里），这个变量就是全局变量
type upvalue struct {
	val *luaValue
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	c := &closure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

func newGoClosure(f GoFunction, nUpvals int) *closure {
	c := &closure{goFunc: f}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}
