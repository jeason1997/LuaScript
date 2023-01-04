/*
 *指令的具体实现：for循环相关
 */
package vm

import . "luago/api"

func forPrep(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	// R(A)-=R(A+2)
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(LUA_OPSUB)
	vm.Replace(a)
	// pc+=sBx
	vm.AddPC(sBx)
}

func forLoop(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	// R(A)+=R(A+2);
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(LUA_OPADD)
	vm.Replace(a)

	// R(A) <? = R(A+1)
	isPositiveStep := vm.ToNumber(a+2) >= 0
	if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
		!isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {

		vm.AddPC(sBx)   // pc+=sBx
		vm.Copy(a, a+3) // R(A+3)=R(A)
	}
}
