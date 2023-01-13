package binchunk

//头部的常量定义
const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

//常量表的TAG定义
const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

//Lua二进制chunk文件的定义
type binaryChunk struct {
	header                  //头部
	sizeUpvalues byte       //主函数upvalue数量
	mainFunc     *Prototype //主函数原型
}

//头部总共占用约30个字节（因平台而异）
type header struct {
	signature       [4]byte //签名，魔数：LUA_SIGNATURE
	version         byte    //Lua版本号，LUAC_VERSION
	format          byte    //格式号，官方的实现格式号，LUAC_FORMAT
	luacData        [6]byte //固定为LUAC_DATA，二次校验用
	cintSize        byte    //cint数据类型在二进制chunk里占用的字节数
	sizetSize       byte    //size_t数据类型在二进制chunk里占用的字节数
	instructionSize byte    //Lua虚拟机指令数在二进制chunk里占用的字节数
	luaIntegerSize  byte    //Lua整数数据类型在二进制chunk里占用的字节数
	luaNumberSize   byte    //Lua浮点数数据类型在二进制chunk里占用的字节数
	luacInt         int64   //存放LUAC_INT，主要用于判断当前机器为大小端
	luacNum         float64 //存放LUAC_NUM，为了检测chunk所使用的浮点数格式
}

//函数原型定义，其中行号表，局部变量表，Upvalue名列表都属于调试信息，可以不需要
//函数原型就相当于面向对象语言里的类，其作用是实例化出真正可执行的函数，也就是闭包。
type Prototype struct {
	Source          string        //源文件名，只有主函数有值
	LineDefined     uint32        //起始行号，主函数为0
	LastLineDefined uint32        //结束行号，主函数为0
	NumParams       byte          //固定参数个数
	IsVararg        byte          //是否为Vararg函数，即是否有变长参数
	MaxStackSize    byte          //寄存器数量
	Code            []uint32      //指令表，每条指令占INSTRUCTION_SIZE个字节
	Constants       []interface{} //常量表，包括nil、布尔值、整数、浮点数和字符串五种
	Upvalues        []Upvalue     //每个占2字节
	Protos          []*Prototype  //子函数原型表
	LineInfo        []uint32      //行号表，与指令表里的指令一一对应
	LocVars         []LocVar      //局部变量表
	UpvalueNames    []string      //Upvalue名列表，与Upvalues表一一对应
}

type Upvalue struct {
	Instack byte //Upvalue捕获的是否是直接外围函数的局部变量，1表示是，0表示否
	Idx     byte //如果Upvalue捕获的是直接外围函数的局部变量，局部变量在外围函数调用帧里的索引
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

//chunk解析函数
func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        // 校验头部
	reader.readByte()           // 跳过Upvalue数量
	return reader.readProto("") // 读取函数原型
}
