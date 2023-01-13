// 测试虚拟机运行Luac文件
package test

import (
	"luago/state"
)

func TestVM(data []byte, chunkName string) {
	ls := state.New()
	ls.Load(data, chunkName, "b")
	//执行主函数
	ls.Call(0, 0)
}
