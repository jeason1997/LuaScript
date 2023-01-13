package vm

import "luago/api"

const MAXARG_Bx = 1<<18 - 1       // 2^18-1 = 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 262143 / 2 = 131071

/*
 31       22       13       5    0
  +-------+^------+-^-----+-^-----
  |b=9bits |c=9bits |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    bx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |   sbx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    ax=26bits            |op=6|
  +-------+^------+-^-----+-^-----
 31      23      15       7      0
*/
type Instruction uint32

//从指令中提取操作码
func (self Instruction) Opcode() int {
	//32位指令中，低6位是操作码，即00000000 00000000 00000000 00xxxxxx
	//跟0x3F(即00111111)做与操作即可提取这部分的值
	return int(self & 0x3F)
}

//从iABC模式指令中提取参数
func (self Instruction) ABC() (a, b, c int) {
	//iABC模式的指令可以携带A、B、C三个操作数，分别占用8、9、9个比特
	a = int(self >> 6 & 0xFF)
	c = int(self >> 14 & 0x1FF)
	b = int(self >> 23 & 0x1FF)
	return
}

//从iABx模式指令中提取参数
func (self Instruction) ABx() (a, bx int) {
	//iABx模式的指令可以携带A和Bx两个操作数，分别占用8和18个比特
	a = int(self >> 6 & 0xFF)
	bx = int(self >> 14)
	return
}

//从iAsBx模式指令中提取参数
func (self Instruction) AsBx() (a, sbx int) {
	//iAsBx模式的指令可以携带A和sBx两个操作数，分别占用8和18个比特
	a, bx := self.ABx()
	return a, bx - MAXARG_sBx
}

//从iAx模式指令中提取参数
func (self Instruction) Ax() int {
	//iAx模式的指令只携带一个操作数，占用全部的26个比特
	return int(self >> 6)
}

func (self Instruction) OpName() string {
	return opcodes[self.Opcode()].name
}

func (self Instruction) OpMode() byte {
	return opcodes[self.Opcode()].opMode
}

func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argBMode
}

func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}

func (self Instruction) Execute(vm api.LuaVM) {
	action := opcodes[self.Opcode()].action
	if action != nil {
		action(self, vm)
	} else {
		panic(self.OpName())
	}
}
