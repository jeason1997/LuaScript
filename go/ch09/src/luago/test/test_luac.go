// 测试加载并解析luac文件
package test

import (
	"luago/binchunk"
	debugger "luago/utils"
)

func TestUndump(data []byte) {
	proto := binchunk.Undump(data)
	list(proto)
}

func list(f *binchunk.Prototype) {
	debugger.PrintHeader(f)
	debugger.PrintCode(f)
	debugger.PrintDetail(f)
	for _, p := range f.Protos {
		list(p)
	}
}
