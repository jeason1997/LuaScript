/*
	 *该脚本是luago/api/lua_vm.go里的接口的具体实现
	 *主要实现：虚拟机的基础操作
	 	PC() int          // 返回当前PC（仅测试用）
		AddPC(n int)      // 修改PC（用于实现跳转指令）
		Fetch() uint32    // 取出当前指令；将PC指向下一条指令
		GetConst(idx int) // 将指定常量推入栈顶
		GetRK(rk int)     // 将指定常量或栈值推入栈顶
*/
package state

func (self *luaState) PC() int {
	return self.stack.pc
}

func (self *luaState) AddPC(n int) {
	self.stack.pc += n
}

// 根据PC索引从函数原型的指令表里取出当前指令，然后把PC加1，这样下次再调用该方法取出的就是下一条指令
func (self *luaState) Fetch() uint32 {
	i := self.stack.closure.proto.Code[self.stack.pc]
	self.stack.pc++
	return i
}

// 根据索引从函数原型的常量表里取出一个常量值，然后把它推入栈顶
func (self *luaState) GetConst(idx int) {
	c := self.stack.closure.proto.Constants[idx]
	self.stack.push(c)
}

//根据情况调用GetConst()方法把某个常量推入栈顶，或者调用PushValue()方法把某个索引处的栈值推入栈顶
func (self *luaState) GetRK(rk int) {
	/*
	 *传递给GetRK()方法的参数实际上是iABC模式指令里的OpArgK类型参数。
	 *由第3章可知，这种类型的参数一共占9个比特。如果最高位是1，那么参数里存放的是常量表索引，把最高位去掉就可以得到索引值；
	 *否则最高位是0，参数里存放的就是寄存器索引值。
	 *但是请读者留意，Lua虚拟机指令操作数里携带的寄存器索引是从0开始的，
	 *而Lua API里的栈索引是从1开始的，所以当需要把寄存器索引当成栈索引使用时，要对寄存器索引加1。
	 */
	if rk > 0xFF {
		// 常量
		self.GetConst(rk & 0xFF)
	} else {
		// 寄存器
		self.PushValue(rk + 1)
	}
}
