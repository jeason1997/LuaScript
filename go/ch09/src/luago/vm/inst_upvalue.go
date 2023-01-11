package vm

import . "luago/api"

//把某个全局变量放入指定寄存器
func getTabUp(i Instruction, vm LuaVM) {
	//a:目标寄存器 c:全局变量在表中的key
	a, _, c := i.ABC()
	a += 1

	//把全局变量表推入栈顶
	vm.PushGlobalTable()
	//获取全局变量的名称并推入栈顶
	vm.GetRK(c)
	//根据栈顶弹出的key获取全局变量表中该变量的值并推入栈顶
	vm.GetTable(-2)
	//弹出值并覆盖寄存器a处的值
	vm.Replace(a)
	//把推入的全局表弹出，恢复栈状态
	vm.Pop(1)
}
