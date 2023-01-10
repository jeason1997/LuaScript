package vm

import . "luago/api"

/*
 *SELF指令主要用来优化方法调用语法糖。
 *比如说obj:f(a, b, c)，虽然从语义的角度来说完全等价于obj.f(obj, a, b, c)，
 *但是Lua编译器并不是先去掉语法糖再按普通的函数调用处理，而是会生成SELF指令，这样就可以节约一条指令。
 *SELF指令（iABC模式）把对象和方法拷贝到相邻的两个目标寄存器中。
 *对象在寄存器中，索引由操作数B指定。方法名在常量表里，索引由操作数C指定。目标寄存器索引由操作数A指定。
 *R(A+1) := R(B); R(A) := R(B)[RK(C)]
 */
func self(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	//把对象从b中拷贝到a+1里
	vm.Copy(b, a+1)
	//获取方法名，并推入栈顶
	vm.GetRK(c)
	//根据栈顶的key（方法名），从表里（即存在b处的对象）获取对应的函数，并且推入栈顶
	vm.GetTable(b)
	//用栈顶的函数覆盖寄存器a处的值
	vm.Replace(a)

	/*
	 *所以最终寄存器的状态就是 [...]->[a:函数闭包]->[a+1:obj]->[...]
	 *obj:f(a, b, c)跟obj.f(obj, a, b, c)对比，最终栈里的样子是完全一模一样的。
	 */
}

/*
 *CLOSURE指令（iBx模式）把当前Lua函数的子函数原型实例化为闭包，放入由操作数A指定的寄存器中。
 *子函数原型来自于当前函数原型的子函数原型表，索引由操作数Bx指定
 *CLOSURE指令对应Lua脚本里的函数定义语句或者表达式
 *R(A) := closure(KPROTO[Bx])
 */
func closure(i Instruction, vm LuaVM) {
	//解码指令：a目标寄存器索引，bx函数原型表索引
	a, bx := i.ABx()
	a += 1

	//加载原型并推入到栈顶
	vm.LoadProto(bx)
	//将栈顶的原型覆盖寄存器a处的值
	vm.Replace(a)
}

/*
 *VARARG指令（iABC模式）把传递给当前函数的变长参数加载到连续多个寄存器中。
 *传递给函数的变长参数一开始是存在调用帧里的，stack.varargs，有需要的时候才突入到栈里
 *其中第一个寄存器的索引由操作数A指定，寄存器数量由操作数B指定
 *R(A), R(A+1), ..., R(A+B-2) = vararg
 */
func vararg(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	//起始寄存器索引A，结束寄存器索引：N=A+B-2
	//如果B=1，则N=A-1，说明不返回参数，如果B=2，则N=A，说明返回一个参数，即A所处位置的值
	if b != 1 {
		//操作数B若大于1，表示把B -1个vararg参数复制到寄存器；否则只能等于0，表示把全部vararg参数复制到寄存器
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

// return R(A)(R(A+1), ... ,R(A+B-1))
func tailCall(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	// todo: optimize tail call!
	c := 0
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

/*
 *CALL指令（iABC模式）调用Lua函数。其中被调函数位于寄存器中，索引由操作数A指定。
 *需要传递给被调函数的参数值也在寄存器中，紧挨着被调函数，数量由操作数B指定。
 *函数调用结束后，原先存放函数和参数值的寄存器会被返回值占据，具体有多少个返回值则由操作数C指定。
 *R(A), ... , R(A+C-2) := R(A)(R(A+1), ... , R(A+B-1))
 */
func call(i Instruction, vm LuaVM) {
	//解码指令：a:被调函数在寄存器的索引，b:需要传入的参数数量，c:返回参数数量
	a, b, c := i.ABC()
	a += 1

	// println(":::"+ vm.StackToString())
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

//a:被调函数在寄存器的索引
//b:需要传入的参数数量
func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	if b >= 1 {
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else {
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

func _fixStack(a int, vm LuaVM) {
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

func _popResults(a, c int, vm LuaVM) {
	if c == 1 {
		// no results
	} else if c > 1 {
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		// leave results on stack
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

/*
 *RETURN指令（iABC模式）把存放在连续多个寄存器里的值返回给主调函数。
 *其中第一个寄存器的索引由操作数A指定，寄存器数量由操作数B指定，操作数C没用
 *return R(A), ... ,R(A+B-2)
 */
func _return(i Instruction, vm LuaVM) {
	//解码指令：a起始寄存器索引，b决定寄存器数量
	a, b, _ := i.ABC()
	a += 1

	//起始寄存器索引A，结束寄存器索引：N=A+B-2
	//如果B=1，则N=A-1，说明不返回参数，如果B=2，则N=A，说明返回一个参数，即A所处位置的值
	if b == 1 {
		// 如果操作数B等于1，N=A+B-2=A-1，则不需要返回任何值
	} else if b > 1 {
		// 如果操作数B大于1，则需要返回B-1个值，这些值已经在寄存器里了，循环调用PushValue()方法复制到栈顶即可
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else {
		//如果操作数B等于0，则一部分返回值已经在栈顶了，调用_fixStack()函数把另一部分也推入栈顶。
		_fixStack(a, vm)
	}
}
