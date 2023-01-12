/*
 *指令的具体实现：其他类指令
 */
package vm

import . "luago/api"

/*
 *MOVE指令（iABC模式）把源寄存器（索引由操作数B指定）里的值移动到目标寄存器（索引由操作数A指定）里
 *MOVE指令常用于局部变量赋值和参数传递
 *由于MOVE等指令使用操作数A（占8个比特）来表示目标寄存器索引，所以Lua函数使用的局部变量不能超过255个
 *opcode{0, 1, OpArgR, OpArgN, IABC, "MOVE     "},
 */
func move(i Instruction, vm LuaVM) {
	//解码指令，得到目标寄存器a和源寄存器索引b
	a, b, _ := i.ABC()
	//寄存器索引（从0开始）加1才是相应的栈索引（从1开始）
	a += 1
	b += 1
	vm.Copy(b, a)
}

/*
 *JMP指令（iAsBx模式）执行无条件跳转。该指令往往和后面要介绍的TEST等指令配合使用，但是也可能会单独出现，比如Lua也支持标签和goto语句
 *opcode{0, 0, OpArgR, OpArgN, IAsBx, "JMP      "},
 */
func jmp(i Instruction, vm LuaVM) {
	//解码指令，得到目标地址偏移值sBx
	a, sBx := i.AsBx()
	//跳转到地址，即新pc=当前pc+偏移sBx
	vm.AddPC(sBx)
	if a != 0 {
		//JMP指令的操作数A和Upvalue有关
		panic("todo! ")
	}
}
