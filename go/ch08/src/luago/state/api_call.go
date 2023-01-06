package state

import (
	"fmt"
	"luago/binchunk"
	"luago/vm"
)

/*
 *如果加载的是二进制chunk，那么只要读取文件、解析主函数原型、实例化为闭包、推入栈顶就可以了；
 *如果加载的是Lua脚本，则要先进行编译（暂时不支持加载lua脚本）
 *chunk：要加载的chunk数据
 *chunkName：指定chunk的名字，供加载错误或调试时使用
 *mode：加载模式（b、t、bt）
 *		-b：第一个参数必须是二进制chunk数据，否则加载失败
 *		-t：第一个参数必须是文本chunk数据，否则加载失败
 *		-bt：第一个参数可以是二进制或者文本chunk数据，会根据实际的数据格式进行处理
 *return：
 *		-0：加载成功
 */
func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	//把主函数原型实例化为闭包并推入栈顶。
	c := newLuaClosure(proto)
	self.stack.push(c)
	return 0
}

/*
 *对Lua函数进行调用。在执行Call()方法之前，必须先把被调函数推入栈顶，然后把参数值依次推入栈顶。
 *Call()方法结束之后，参数值和函数会被弹出栈顶，取而代之的是指定数量的返回值压入栈顶。
 *nArgs：准备传递给被调函数的参数数量，同时也隐含给出了被调函数在栈里的位置
 *nResults：需要的返回值数量（多退少补），如果是-1，则被调函数的返回值会全部留在栈顶。
 */
func (self *luaState) Call(nArgs, nResults int) {
	//此时栈里的状态是，传参在栈顶，接下来是被调函数，因此可以通过栈顶减去参数的数量来获得被调函数的位置
	val := self.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		fmt.Printf("call %s<%d,%d>\n", c.proto.Source,
			c.proto.LineDefined, c.proto.LastLineDefined)
		self.callLuaClosure(nArgs, nResults, c)
	} else {
		panic("not function!")
	}
}

func (self *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	//从函数原型里获取各种信息：执行函数所需的寄存器数量，固定参数数量，是否有不定参
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	//根据寄存器数量（适当扩大，因为要给指令实现函数预留少量栈空间）创建一个新的调用帧，并把闭包和调用帧联系起来
	newStack := newLuaStack(nRegs + 20)
	newStack.closure = c

	//把函数和参数值一次性从栈顶弹出，然后调用新帧的pushN()方法按照固定参数数量传入参数
	funcAndArgs := self.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	//固定参数传递完毕之后，需要修改新帧的栈顶指针，让它指向最后一个寄存器
	newStack.top = nRegs
	//如果被调函数是vararg函数，且传入参数的数量多于固定参数数量，还需要把vararg参数记下来，存在调用帧里，以备后用
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	//把新调用帧推入调用栈顶，让它成为当前帧
	self.pushLuaStack(newStack)
	//执行被调函数的指令
	self.runLuaClosure()
	//指令执行完毕之后，新调用帧的使命就结束了，把它从调用栈顶弹出，这样主调帧就又成了当前帧
	self.popLuaStack()

	//被调函数运行完毕之后，返回值会留在被调帧的栈顶（寄存器之上）。
	//我们需要把全部返回值从被调帧栈顶弹出，然后根据期望的返回值数量多退少补，推入当前帧栈顶
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		//不够放参数的话就扩容
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

func (self *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(self.Fetch())
		inst.Execute(self)

		//打印调试信息
		pc := self.PC()
		fmt.Printf("[%02d] %s ", pc+1, inst.OpName())
		//printStack(self)

		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}
