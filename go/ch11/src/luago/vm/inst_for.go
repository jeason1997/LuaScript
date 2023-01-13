/*
 *指令的具体实现：for循环相关
 *Lua语言的for循环语句有两种形式：数值（Numerical）形式和通用（Generic）形式。
 *数值for循环用于按一定步长遍历某个范围内的数值，通用for循环主要用于遍历表。
 *数值for循环需要借助两条指令来实现：FORPREP和FORLOOP
 */
package vm

import . "luago/api"

//FORPREP指令执行的操作其实就是在循环开始之前预先给数值减去步长，然后跳转到FORLOOP指令正式开始循环
func forPrep(i Instruction, vm LuaVM) {
	/*
	 * for i=index, limit, setp do f() end
	 * a = index
	 * a + 1 = limit
	 * a + 2 = step
	 * a + 3 = i
	 */
	//解码指令，a为存放index值的寄存器索引，sBx为FORLOOP处处的指令偏移地址
	a, sBx := i.AsBx()
	a += 1

	// index减去步长，R(A)-=R(A+2)
	vm.PushValue(a)     //index
	vm.PushValue(a + 2) //step
	vm.Arith(LUA_OPSUB) //index-step
	vm.Replace(a)
	// 跳转到FORLOOP处开始循环，pc+=sBx
	vm.AddPC(sBx)
}

//FORLOOP指令则是先给数值加上步长，然后判断数值是否还在范围之内。
//如果已经超出范围，则循环结束；若未超过范围则把数值拷贝给用户定义的局部变量，然后跳转到循环体内部开始执行具体的代码块
func forLoop(i Instruction, vm LuaVM) {
	/*
	 * for i=index, limit, setp do f() end
	 * a = index
	 * a + 1 = limit
	 * a + 2 = step
	 * a + 3 = i
	 */
	//解码指令，a为存放index值的寄存器索引，sBx为循环开始处的指令偏移地址
	a, sBx := i.AsBx()
	a += 1

	// index加上步长，R(A)+=R(A+2);
	vm.PushValue(a + 2) //setp
	vm.PushValue(a)     //index
	vm.Arith(LUA_OPADD) //index+setp
	vm.Replace(a)

	// R(A) <? = R(A+1)
	//判断step是正向增长还是反向增长
	isPositiveStep := vm.ToNumber(a+2) >= 0

	//step是正向增长并且index小于等于limit
	//或者step是反向增长并且index大于limit
	if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
		!isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {
		//下一条指令跳转到循环的开头
		vm.AddPC(sBx) // pc+=sBx
		//i赋值为index
		vm.Copy(a, a+3) // R(A+3)=R(A)
	} else {
		//否则直接结束指令，执行下一条指令，即结束循环
	}
}
