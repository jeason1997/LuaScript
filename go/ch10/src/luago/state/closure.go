package state

import (
	. "luago/api"
	"luago/binchunk"
)

// 我们使用closure结构体来统一表示Lua和Go闭包。
// 如果proto字段不是nil，说明这是Lua闭包。否则，goFunc字段一定不是nil，说明这是Go闭包。
type closure struct {
	proto  *binchunk.Prototype //Lua闭包
	goFunc GoFunction          //Go闭包
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	return &closure{proto: proto}
}

func newGoClosure(f GoFunction) *closure {
	return &closure{goFunc: f}
}
