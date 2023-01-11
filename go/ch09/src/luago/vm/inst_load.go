/*
 *指令的具体实现：加载类指令
 *加载指令用于把nil值、布尔值或者常量表里的常量值加载到寄存器里。
 */
package vm

import . "luago/api"

/*
 *LOADNIL指令（i ABC模式）用于给连续n个寄存器放置nil值。寄存器的起始索引由操作数A指定，寄存器数量则由操作数B指定，操作数C没有用
 *在Lua代码里，局部变量的默认初始值是nil。LOADNIL指令常用于给连续n个局部变量设置初始值
 *opcode{0, 1, OpArgU, OpArgN, IABC, "LOADNIL "}
 */
func loadNil(i Instruction, vm LuaVM) {
	//解码指令
	a, b, _ := i.ABC()
	//寄存器索引（从0开始）加1才是相应的栈索引（从1开始）
	a += 1

	//假定虚拟机在执行第一条指令前，已经预先算好执行阶段所需要的寄存器数量，调用SetTop()方法保留了必要数量的栈空间
	//先调用PushNil()方法往栈顶推入一个nil值，然后连续调用Copy()方法将该nil值复制到指定寄存器中，
	//最后调用Pop()方法把一开始推入栈顶的那个nil值弹出，让栈顶指针恢复原状
	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

/*
 *LOADBOOL指令（iABC模式）给单个寄存器设置布尔值。
 *寄存器索引由操作数A指定，布尔值由寄存器B指定（0代表false，非0代表true），如果寄存器C非0则跳过下一条指令。
 *opcode{0, 1, OpArgU, OpArgU, IABC, "LOADBOOL"}
 */
func loadBool(i Instruction, vm LuaVM) {
	//解码指令
	a, b, c := i.ABC()
	//寄存器索引（从0开始）加1才是相应的栈索引（从1开始）
	a += 1
	//将布尔值压入栈顶
	vm.PushBoolean(b != 0)
	//用栈顶的布尔值弹出，并覆盖索引a处的值
	vm.Replace(a)
	//如果c非0则跳过下一条指令
	if c != 0 {
		vm.AddPC(1)
	}
}

/*
 *LOADK指令（iABx模式）将常量表里的某个常量加载到指定寄存器，寄存器索引由操作数A指定，常量表索引由操作数Bx指定。
 *opcode{0, 1, OpArgK, OpArgN, IABx, "LOADK    "}
 */
func loadK(i Instruction, vm LuaVM) {
	//解码指令，得到目标寄存器a和常量索引bx
	a, bx := i.ABx()
	//寄存器索引（从0开始）加1才是相应的栈索引（从1开始）
	a += 1
	//将常量表里bx索引处的常量推向栈顶
	vm.GetConst(bx)
	//用栈顶的常量覆盖索引a处的值
	vm.Replace(a)
}

/*
 *LOADK指令，操作数Bx占18个比特，能表示的最大无符号整数是262143，
 *大部分Lua函数的常量表大小都不会超过这个数，所以这个限制通常不是什么问题。
 *不过Lua也经常被当作数据描述语言使用，所以常量表大小可能超过这个限制也并不稀奇。为了应对这种情况，Lua还提供了一条LOADKX指令
 *LOADKX指令（也是iABx模式）需要和EXTRAARG指令（iAx模式）搭配使用，用后者的Ax操作数来指定常量索引。Ax操作数占26个比特，可以表达的最大无符号整数是67108864
 *opcode{0, 1, OpArgN, OpArgN, IABx, "LOADKX  "}
 */
func loadKx(i Instruction, vm LuaVM) {
	a, _ := i.ABx()
	a += 1
	//LOADK指令把操作数放在指令的Bx里，而LOADKX则把操作数放在下一条指令里
	ax := Instruction(vm.Fetch()).Ax()
	vm.GetConst(ax)
	vm.Replace(a)
}
