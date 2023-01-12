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

//返回当前Lua函数所操作的寄存器数量（编译时自动计算最大需要数量）
func (self *luaState) RegisterCount() int {
	return int(self.stack.closure.proto.MaxStackSize)
}

//把传递给当前Lua函数的变长参数推入栈顶（多退少补）
func (self *luaState) LoadVararg(n int) {
	if n < 0 {
		//小于0则将参数全部推入
		n = len(self.stack.varargs)
	}
	self.stack.check(n)
	self.stack.pushN(self.stack.varargs, n)
}

//把当前Lua函数的子函数的原型实例化为闭包推入栈顶
func (self *luaState) LoadProto(idx int) {
	stack := self.stack
	subProto := stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)
	stack.push(closure)

	//加载子函数原型时也需要初始化Upvalue（主函数Main的Upvalue在Load程序的时候已经加载了）
	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 {
			//如果某一个Upvalue捕获的是当前函数的局部变量
			//那么我们只要访问当前函数的局部变量即可

			if stack.openuvs == nil {
				stack.openuvs = map[int]*upvalue{}
			}
			if openuv, found := stack.openuvs[uvIdx]; found {
				//如果Upvalue捕获的外围函数局部变量还在栈上，直接引用即可，我们称这种Upvalue处于开放（Open）状态
				closure.upvals[i] = openuv
			} else {
				//反之，必须把变量的实际值保存在其他地方，我们称这种Upvalue处于闭合（Closed）状态
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				//为了能够在合适的时机（比如局部变量退出作用域时）把处于开放状态的Upvalue闭合，
				//需要记录所有暂时还处于开放状态的Upvalue，我们把这些Upvalue记录在被捕获局部变量所在的栈帧里
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
		} else {
			//如果某一个Upvalue捕获的是更外围的函数中的局部变量
			//该Upvalue已经被当前函数捕获，我们只要把该Upvalue传递给闭包即可
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}
}

func (self *luaState) CloseUpvalues(a int) {
	for i, openuv := range self.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(self.stack.openuvs, i)
		}
	}
}
