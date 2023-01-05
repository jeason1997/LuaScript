// 测试虚拟机运行Luac文件
package test

import (
	"fmt"
	"luago/binchunk"
	"luago/state"
	. "luago/vm"
)

func TestVM(data []byte) {
	//加载预编译的二进制chunk文件，并解析出函数原型并执行其中的指令
	proto := binchunk.Undump(data)
	//编译的时候会自动计算运行需要的最大寄存器数量
	nRegs := int(proto.MaxStackSize)
	//创建虚拟机，由于指令实现函数也需要少量的栈空间，所以实际创建的Lua栈容量要比寄存器数量稍微大一些
	ls := state.New(nRegs+8, proto)
	//调用SetTop()方法在栈里预留出寄存器空间，剩余栈空间留给指令实现函数使用
	ls.SetTop(nRegs)
	/*
	*一个运行于虚拟机的正常栈，里面的结构应该是，底部是程序预留的寄存器（编译时自动计算最多需要多少个寄存器Prototype.MaxStackSize）。
	*上面剩余的是计算用到的栈空间，一般会预留几个。然后初始时栈顶索引是位于寄存器上面的，也就是计算栈的初始位置。
	*slots = [reg1][reg2][reg3][stack1][stack2][...]
	*top = 4
	 */

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
