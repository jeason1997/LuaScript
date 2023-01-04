/*
 *指令的具体实现：运算符相关
 */
package vm

import . "luago/api"

/*
 *二元算术运算指令（iABC模式），对两个寄存器或常量值（索引由操作数B和C指定）进行运算，将结果放入另一个寄存器（索引由操作数A指定）
	LUA_OPADD  = iota // +
	LUA_OPSUB         // -
	LUA_OPMUL         // ＊
	LUA_OPMOD         // %
	LUA_OPPOW         // ^
	LUA_OPDIV         // /
	LUA_OPIDIV        // //
	LUA_OPBAND        // &
	LUA_OPBOR         // |
	LUA_OPBXOR        // ～
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
*/
func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
	//解码指令
	a, b, c := i.ABC()
	//寄存器索引（从0开始）加1才是相应的栈索引（从1开始）
	a += 1
	//将两个操作数里的常量或者寄存器值推入栈顶
	vm.GetRK(b)
	vm.GetRK(c)
	//进行算术运算，并将结果推入栈顶
	vm.Arith(op)
	//将栈顶的值覆盖栈里索引a处的值
	vm.Replace(a)
}

/*
 *一元算术运算指令（iABC模式），对操作数B所指定的寄存器里的值进行运算，然后把结果放入操作数A所指定的寄存器中，操作数C没用
 *为什么二元运算的两个操作数可以是寄存器也可以是常量，而一元运算的操作数只能是寄存器呢？
 *因为一元运算常量没有意义，在编译阶段可以直接将一元运算的常量值计算出来，比如-1，直接当成负1的数值，没必要进行计算
	LUA_OPUNM         // - (一元取反)
	LUA_OPBNOT        // ～
*/
func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.PushValue(b)
	vm.Arith(op)
	vm.Replace(a)
}

/*
 *比较指令（iABC模式），比较寄存器或常量表里的两个值（索引分别由操作数B和C指定），如果比较结果和操作数A（转换为布尔值）匹配，则跳过下一条指令。比较指令不改变寄存器状态
 *比较指令对应Lua语言里的比较运算符（当用于赋值时，需要和LOADBOOL指令搭配使用）
 */
func _compare(i Instruction, vm LuaVM, op CompareOp) {
	a, b, c := i.ABC()
	//将两个操作数里的常量或者寄存器值推入栈顶
	vm.GetRK(b)
	vm.GetRK(c)
	//比较栈顶的两个值，如果跟操作数a不一样，则跳过下一条指令
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}
	//比较指令不改变寄存器状态，因为上面为了作比较往栈顶推入了2个值，这里需要回退两个值
	vm.Pop(2)
}

func add(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPADD) }  // +
func sub(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSUB) }  // -
func mul(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMUL) }  // ＊
func mod(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMOD) }  // %
func pow(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPPOW) }  // ^
func div(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPDIV) }  // /
func idiv(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPIDIV) } // //
func band(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBAND) } // &
func bor(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPBOR) }  // |
func bxor(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBXOR) } // ～
func shl(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHL) }  // <<
func shr(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHR) }  // >>
func unm(i Instruction, vm LuaVM)  { _unaryArith(i, vm, LUA_OPUNM) }   // -
func bnot(i Instruction, vm LuaVM) { _unaryArith(i, vm, LUA_OPBNOT) }  // ～

func eq(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPEQ) } // ==
func lt(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLT) } // <
func le(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLE) } // <=

/*
 *LEN指令（iABC模式）进行的操作和一元算术运算指令类似
 *同样操作数只能是寄存器，常量没意义，因为可以在编译阶段直接算出来，例如#'abc'，直接用常量3代替
 *opcode{0, 1, OpArgR, OpArgN, IABC, "LEN      "}
 */
func length(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Len(b)
	vm.Replace(a)
}

/*
 *CONCAT指令（iABC模式），将连续n个寄存器（起止索引分别由操作数B和C指定）里的值拼接，将结果放入另一个寄存器（索引由操作数A指定）
 *opcode{0, 1, OpArgR, OpArgR, IABC, "CONCAT  "}
 */
func concat(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	c += 1
	//数量
	n := c - b + 1
	//检查栈空间是否足够，不够的话扩容
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		//将要拼接的值，依次放到栈顶
		vm.PushValue(i)
	}
	//从栈顶弹出n个值，对这些值进行拼接，然后把结果推入栈顶
	vm.Concat(n)
	//将栈顶的结果覆盖到栈索引a处的值
	vm.Replace(a)
}

/*
 *NOT指令（iABC模式）进行的操作和一元算术运算指令类似
 *对操作数B所指定的寄存器里的值进行取反，然后把结果放入操作数A所指定的寄存器中\
 *opcode{0, 1, OpArgR, OpArgN, IABC, "NOT      "}
 */
func not(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

/*
 *TESTSET指令（iABC模式），判断寄存器B（索引由操作数B指定）中的值转换为布尔值之后是否和操作数C表示的布尔值一致，如果一致则将寄存器B中的值复制到寄存器A（索引由操作数A指定）中，否则跳过下一条指令。
 *TESTSET指令对应Lua语言里的逻辑与和逻辑或运算符，比如 a = b and c
 *opcode{1, 1, OpArgR, OpArgU, IABC, "TESTSET "}
 */
func testSet(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	//如果b等于c，则将b拷贝到a中，否则跳过下一条指令
	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

/*
 *TEST指令（iABC模式），判断寄存器A（索引由操作数A指定）中的值转换为布尔值之后是否和操作数C表示的布尔值一致，如果一致，则跳过下一条指令。TEST指令不使用操作数B，也不改变寄存器状态
 *TEST指令是TESTSET指令的特殊形式，比如 b = b and c
 *opcode{1, 0, OpArgN, OpArgU, IABC, "TEST     "}
 */
func test(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a += 1
	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}
