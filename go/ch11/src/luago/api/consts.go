package api

//LUA数据类型
const (
	LUA_TNONE = iota - 1 // -1
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

//LUA运算类型
const (
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
	LUA_OPUNM         // - (一元取反)
	LUA_OPBNOT        // ～
)

//LUA比较类型（!=、>、>=可以通过==、<、<=拓展得到）
const (
	LUA_OPEQ = iota // ==
	LUA_OPLT        // <
	LUA_OPLE        // <=
)

const LUA_MINSTACK = 20                         //LUA调用栈最小容量
const LUAI_MAXSTACK = 1000000                   //LUA调用栈最大容量（可正负）
const LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000 //负有效索引减1000就是注册表的伪索引
const LUA_RIDX_GLOBALS int64 = 2                //全局环境在注册表里的索引
