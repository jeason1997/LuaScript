// 测试虚拟机
package test

import (
	"fmt"
	"luago/binchunk"
	"luago/state"
	. "luago/vm"
)

func TestVM(data []byte) {
	proto := binchunk.Undump(data)
	//编译的时候会自动计算运行需要的最大寄存器数量
	nRegs := int(proto.MaxStackSize)
	//创建虚拟机，由于指令实现函数也需要少量的栈空间，所以实际创建的Lua栈容量要比寄存器数量稍微大一些
	ls := state.New(nRegs+8, proto)
	//调用SetTop()方法在栈里预留出寄存器空间，剩余栈空间留给指令实现函数使用
	ls.SetTop(nRegs)

	for {
		pc := ls.PC()
		//取一条指令，并递增PC
		inst := Instruction(ls.Fetch())
		if inst.Opcode() != OP_RETURN {
			//执行指令
			inst.Execute(ls)
			fmt.Printf("[%02d] %s ", pc+1, inst.OpName())
			printStack(ls)
		} else {
			break
		}
	}
}
