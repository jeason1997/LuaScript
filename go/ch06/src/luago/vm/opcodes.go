package vm

import "luago/api"

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
	action   func(i Instruction, vm api.LuaVM)
}

//指令集合，与操作码一一对应
var opcodes = []opcode{
	/*     T  A    B       C     mode         name       action */
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "MOVE    ", move},     // R(A) := R(B)
	opcode{0, 1, OpArgK, OpArgN, IABx /* */, "LOADK   ", loadK},    // R(A) := Kst(Bx)
	opcode{0, 1, OpArgN, OpArgN, IABx /* */, "LOADKX  ", loadKx},   // R(A) := Kst(extra arg)
	opcode{0, 1, OpArgU, OpArgU, IABC /* */, "LOADBOOL", loadBool}, // R(A) := (bool)B; if (C) pc++
	opcode{0, 1, OpArgU, OpArgN, IABC /* */, "LOADNIL ", loadNil},  // R(A), R(A+1), ..., R(A+B) := nil
	opcode{0, 1, OpArgU, OpArgN, IABC /* */, "GETUPVAL", nil},      // R(A) := UpValue[B]
	opcode{0, 1, OpArgU, OpArgK, IABC /* */, "GETTABUP", nil},      // R(A) := UpValue[B][RK(C)]
	opcode{0, 1, OpArgR, OpArgK, IABC /* */, "GETTABLE", nil},      // R(A) := R(B)[RK(C)]
	opcode{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABUP", nil},      // UpValue[A][RK(B)] := RK(C)
	opcode{0, 0, OpArgU, OpArgN, IABC /* */, "SETUPVAL", nil},      // UpValue[B] := R(A)
	opcode{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABLE", nil},      // R(A)[RK(B)] := RK(C)
	opcode{0, 1, OpArgU, OpArgU, IABC /* */, "NEWTABLE", nil},      // R(A) := {} (size = B,C)
	opcode{0, 1, OpArgR, OpArgK, IABC /* */, "SELF    ", nil},      // R(A+1) := R(B); R(A) := R(B)[RK(C)]
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "ADD     ", add},      // R(A) := RK(B) + RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SUB     ", sub},      // R(A) := RK(B) - RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "MUL     ", mul},      // R(A) := RK(B) * RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "MOD     ", mod},      // R(A) := RK(B) % RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "POW     ", pow},      // R(A) := RK(B) ^ RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "DIV     ", div},      // R(A) := RK(B) / RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "IDIV    ", idiv},     // R(A) := RK(B) // RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BAND    ", band},     // R(A) := RK(B) & RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BOR     ", bor},      // R(A) := RK(B) | RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BXOR    ", bxor},     // R(A) := RK(B) ~ RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SHL     ", shl},      // R(A) := RK(B) << RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SHR     ", shr},      // R(A) := RK(B) >> RK(C)
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "UNM     ", unm},      // R(A) := -R(B)
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "BNOT    ", bnot},     // R(A) := ~R(B)
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "NOT     ", not},      // R(A) := not R(B)
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "LEN     ", length},   // R(A) := length of R(B)
	opcode{0, 1, OpArgR, OpArgR, IABC /* */, "CONCAT  ", concat},   // R(A) := R(B).. ... ..R(C)
	opcode{0, 0, OpArgR, OpArgN, IAsBx /**/, "JMP     ", jmp},      // pc+=sBx; if (A) close all upvalues >= R(A - 1)
	opcode{1, 0, OpArgK, OpArgK, IABC /* */, "EQ      ", eq},       // if ((RK(B) == RK(C)) ~= A) then pc++
	opcode{1, 0, OpArgK, OpArgK, IABC /* */, "LT      ", lt},       // if ((RK(B) <  RK(C)) ~= A) then pc++
	opcode{1, 0, OpArgK, OpArgK, IABC /* */, "LE      ", le},       // if ((RK(B) <= RK(C)) ~= A) then pc++
	opcode{1, 0, OpArgN, OpArgU, IABC /* */, "TEST    ", test},     // if not (R(A) <=> C) then pc++
	opcode{1, 1, OpArgR, OpArgU, IABC /* */, "TESTSET ", testSet},  // if (R(B) <=> C) then R(A) := R(B) else pc++
	opcode{0, 1, OpArgU, OpArgU, IABC /* */, "CALL    ", nil},      // R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
	opcode{0, 1, OpArgU, OpArgU, IABC /* */, "TAILCALL", nil},      // return R(A)(R(A+1), ... ,R(A+B-1))
	opcode{0, 0, OpArgU, OpArgN, IABC /* */, "RETURN  ", nil},      // return R(A), ... ,R(A+B-2)
	opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORLOOP ", forLoop},  // R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }
	opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORPREP ", forPrep},  // R(A)-=R(A+2); pc+=sBx
	opcode{0, 0, OpArgN, OpArgU, IABC /* */, "TFORCALL", nil},      // R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
	opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "TFORLOOP", nil},      // if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx }
	opcode{0, 0, OpArgU, OpArgU, IABC /* */, "SETLIST ", nil},      // R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
	opcode{0, 1, OpArgU, OpArgN, IABx /* */, "CLOSURE ", nil},      // R(A) := closure(KPROTO[Bx])
	opcode{0, 1, OpArgU, OpArgN, IABC /* */, "VARARG  ", nil},      // R(A), R(A+1), ..., R(A+B-2) = vararg
	opcode{0, 0, OpArgU, OpArgU, IAx /*  */, "EXTRAARG", nil},      // extra (larger) argument for previous opcode
}
