package vm

//编码模式
const (
	IABC  = iota //iABC模式的指令可以携带A、B、C三个操作数，分别占用8、9、9个比特
	IABx         //iABx模式的指令可以携带A和Bx两个操作数，分别占用8和18个比特
	IAsBx        //iAsBx模式的指令可以携带A和sBx两个操作数，分别占用8和18个比特
	IAx          //iAx模式的指令只携带一个操作数，占用全部的26个比特
)

//操作码
const (
	OP_MOVE = iota
	OP_LOADK
	OP_LOADKX
	OP_LOADBOOL
	OP_LOADNIL
	OP_GETUPVAL
	OP_GETTABUP
	OP_GETTABLE
	OP_SETTABUP
	OP_SETUPVAL
	OP_SETTABLE
	OP_NEWTABLE
	OP_SELF
	OP_ADD
	OP_SUB
	OP_MUL
	OP_MOD
	OP_POW
	OP_DIV
	OP_IDIV
	OP_BAND
	OP_BOR
	OP_BXOR
	OP_SHL
	OP_SHR
	OP_UNM
	OP_BNOT
	OP_NOT
	OP_LEN
	OP_CONCAT
	OP_JMP
	OP_EQ
	OP_LT
	OP_LE
	OP_TEST
	OP_TESTSET
	OP_CALL
	OP_TAILCALL
	OP_RETURN
	OP_FORLOOP
	OP_FORPREP
	OP_TFORCALL
	OP_TFORLOOP
	OP_SETLIST
	OP_CLOSURE
	OP_VARARG
	OP_EXTRAARG
)

//操作数类型
const (
	OpArgN = iota // OpArgN类型的操作数不表示任何信息，也就是说不会被使用
	OpArgU        // OpArgU可能表示布尔值、整数值、upvalue索引、子函数索引等
	OpArgR        // OpArgR类型的操作数在iABC模式下表示寄存器索引，在iAsBx模式下表示跳转偏移
	OpArgK        // OpArgK类型的操作数表示常量表索引或者寄存器索引，具体可以分为两种情况。第一种情况是LOADK指令（iABx模式，用于将常量表中的常量加载到寄存器中），该指令的Bx操作数表示常量表索引；第二种情况是部分iABC模式指令，这些指令的B或C操作数既可以表示常量表索引也可以表示寄存器索引，通过最高位是否为1来表示是否常量索引
)

//指令结构
type opcode struct {
	testFlag byte // operator is a test (next instruction must be a jump)
	setAFlag byte // instruction set register A
	argBMode byte // B arg mode
	argCMode byte // C arg mode
	opMode   byte // op mode
	name     string
}

//指令集合，与操作码一一对应
var opcodes = []opcode{
	/*     T  A  B    	 C       mode  name     */
	opcode{0, 1, OpArgR, OpArgN, IABC, "MOVE     "},
	opcode{0, 1, OpArgK, OpArgN, IABx, "LOADK    "},
	opcode{0, 1, OpArgN, OpArgN, IABx, "LOADKX  "},
	opcode{0, 1, OpArgU, OpArgU, IABC, "LOADBOOL"},
	opcode{0, 1, OpArgU, OpArgN, IABC, "LOADNIL "},
	opcode{0, 1, OpArgU, OpArgN, IABC, "GETUPVAL"},
	opcode{0, 1, OpArgU, OpArgK, IABC, "GETTABUP"},
	opcode{0, 1, OpArgR, OpArgK, IABC, "GETTABLE"},
	opcode{0, 0, OpArgK, OpArgK, IABC, "SETTABUP"},
	opcode{0, 0, OpArgU, OpArgN, IABC, "SETUPVAL"},
	opcode{0, 0, OpArgK, OpArgK, IABC, "SETTABLE"},
	opcode{0, 1, OpArgU, OpArgU, IABC, "NEWTABLE"},
	opcode{0, 1, OpArgR, OpArgK, IABC, "SELF     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "ADD      "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "SUB      "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "MUL      "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "MOD      "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "POW      "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "DIV      "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "IDIV     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "BAND     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "BOR      "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "BXOR     "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "SHL      "},
	opcode{0, 1, OpArgK, OpArgK, IABC, "SHR      "},
	opcode{0, 1, OpArgR, OpArgN, IABC, "UNM      "},
	opcode{0, 1, OpArgR, OpArgN, IABC, "BNOT     "},
	opcode{0, 1, OpArgR, OpArgN, IABC, "NOT      "},
	opcode{0, 1, OpArgR, OpArgN, IABC, "LEN      "},
	opcode{0, 1, OpArgR, OpArgR, IABC, "CONCAT  "},
	opcode{0, 0, OpArgR, OpArgN, IAsBx, "JMP      "},
	opcode{1, 0, OpArgK, OpArgK, IABC, "EQ       "},
	opcode{1, 0, OpArgK, OpArgK, IABC, "LT       "},
	opcode{1, 0, OpArgK, OpArgK, IABC, "LE       "},
	opcode{1, 0, OpArgN, OpArgU, IABC, "TEST     "},
	opcode{1, 1, OpArgR, OpArgU, IABC, "TESTSET "},
	opcode{0, 1, OpArgU, OpArgU, IABC, "CALL     "},
	opcode{0, 1, OpArgU, OpArgU, IABC, "TAILCALL"},
	opcode{0, 0, OpArgU, OpArgN, IABC, "RETURN  "},
	opcode{0, 1, OpArgR, OpArgN, IAsBx, "FORLOOP "},
	opcode{0, 1, OpArgR, OpArgN, IAsBx, "FORPREP "},
	opcode{0, 0, OpArgN, OpArgU, IABC, "TFORCALL"},
	opcode{0, 1, OpArgR, OpArgN, IAsBx, "TFORLOOP"},
	opcode{0, 0, OpArgU, OpArgU, IABC, "SETLIST "},
	opcode{0, 1, OpArgU, OpArgN, IABx, "CLOSURE "},
	opcode{0, 1, OpArgU, OpArgN, IABC, "VARARG  "},
	opcode{0, 0, OpArgU, OpArgU, IAx, "EXTRAARG"},
}
