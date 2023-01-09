/*
 *指令的具体实现：table相关
 */
package vm

import . "luago/api"

const LFIELDS_PER_FLUSH = 50

/*
 *NEWTABLE指令（iABC模式）创建空表，并将其放入指定寄存器。
 *寄存器索引由操作数A指定，表的初始数组容量和哈希表容量分别由操作数B和C指定
 *R(A) := {} (size = B, C)
 */
func newTable(i Instruction, vm LuaVM) {
	//解码指令：a目标寄存器索引，b表的数组容量，c表的哈希容量
	a, b, c := i.ABC()
	a += 1
	//创建一个表并推入栈顶
	vm.CreateTable(Fb2int(b), Fb2int(c))
	//将栈顶的表覆盖到寄存器a处
	vm.Replace(a)
}

/*
 *GETTABLE指令（iABC模式）根据键从表里取值，并放入目标寄存器中。
 *其中表位于寄存器中，索引由操作数B指定；键可能位于寄存器中，也可能在常量表里，索引由操作数C指定；
 *目标寄存器索引则由操作数A指定。
 *R(A) := R(B)[RK(C)]
 */
func getTable(i Instruction, vm LuaVM) {
	//解码指令：a目标寄存器索引，b表的寄存器索引，c表的键（可能位于寄存器，也可能位于常量表）
	a, b, c := i.ABC()
	a += 1
	b += 1
	//获取键（常量或者寄存器），并推入到栈顶
	vm.GetRK(c)
	//根据栈顶的key获取表的某个值，并推入栈顶
	vm.GetTable(b)
	//将栈顶的值覆盖寄存器a处的值
	vm.Replace(a)
}

/*
 *SETTABLE指令（iABC模式）根据键往表里赋值。
 *其中表位于寄存器中，索引由操作数A指定；
 *键和值可能位于寄存器中，也可能在常量表里，索引分别由操作数B和C指定。
 *R(A)[RK(B)] := RK(C)
 */
func setTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

/*
 *SETLIST指令（iABC模式）则是专门给数组准备的，用于按索引批量设置数组元素。
 *其中数组位于寄存器中，索引由操作数A指定；需要写入数组的一系列值也在寄存器中，紧挨着数组，数量由操作数B指定；
 *数组起始索引则由操作数C指定。
 *R(A)[(C-1)＊FPF+i] := R(A+i), 1 <= i <= B
 */
func setList(i Instruction, vm LuaVM) {
	//解码指令：a是table的寄存器索引，b需要写入tabel数组的值数量，c要写入的table数组的起始索引
	a, b, c := i.ABC()
	a += 1

	if c > 0 {
		//Lua的table是从1开始，需要减一
		c = c - 1
	} else {
		//如果数组长度大于25600，这种情况下SETLIST指令后面会跟一条EXTRAARG指令，用其Ax操作数来保存批次数。
		c = Instruction(vm.Fetch()).Ax()
	}

	bIsZero := b == 0
	if bIsZero {
		b = int(vm.ToInteger((-1))) - a - 1
		vm.Pop(1)
	}

	//数组的起始位置：因为C操作数只有9个比特（512），所以直接用它表示数组索引显然不够用。
	//这里的解决办法是让C操作数保存批次数，然后用批次数乘上批大小（对应伪代码中的FPF）就可以算出数组起始索引。
	//以默认的批大小50为例，C操作数能表示的最大索引就是扩大到了25600（50*512）。
	vm.CheckStack(1)
	idx := int64(c * LFIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		idx++
		//要连续写入的数据紧靠着table寄存器的位置，依次压入到栈顶
		vm.PushValue(a + j)
		//将栈顶的值写入到table的idx位置
		vm.SetI(a, idx)
	}

	if bIsZero {
		for j := vm.RegisterCount() + 1; j <= vm.GetTop(); j++ {
			idx++
			vm.PushValue(j)
			vm.SetI(a, idx)
		}

		// clear stack
		vm.SetTop(vm.RegisterCount())
	}
}
