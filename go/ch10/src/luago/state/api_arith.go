/*
 *该脚本是luago/api/lua_state.go里的接口的具体实现
 *主要实现：运算操作
 *	Arith(op ArithOp)
 */
package state

import (
	. "luago/api"
	"luago/number"
	"math"
)

var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = number.IMod
	fmod  = number.FMod
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = number.IFloorDiv
	fidiv = number.FFloorDiv
	band  = func(a, b int64) int64 { return a & b }
	bor   = func(a, b int64) int64 { return a | b }
	bxor  = func(a, b int64) int64 { return a ^ b }
	shl   = number.ShiftLeft
	shr   = number.ShiftRight
	iunm  = func(a, _ int64) int64 { return -a }
	funm  = func(a, _ float64) float64 { return -a }
	bnot  = func(a, _ int64) int64 { return ^a }
)

type operator struct {
	integerFunc func(int64, int64) int64       //整数运算
	floatFunc   func(float64, float64) float64 //浮点数运算
}

// 与consts.go里定义的LUA运算类型一一对应
var operators = []operator{
	operator{iadd, fadd},   // +
	operator{isub, fsub},   // -
	operator{imul, fmul},   // ＊
	operator{imod, fmod},   // %
	operator{nil, pow},     // ^
	operator{nil, div},     // /
	operator{iidiv, fidiv}, // //
	operator{band, nil},    // &
	operator{bor, nil},     // |
	operator{bxor, nil},    // ～
	operator{shl, nil},     // <<
	operator{shr, nil},     // >>
	operator{iunm, funm},   // - (一元取反)
	operator{bnot, nil},    // ～
}

// 从栈顶取出操作数，按照一定规则运算，并将结果压回栈顶
func (self *luaState) Arith(op ArithOp) {
	var a, b luaValue

	//不管任何运算，都至少需要一个操作数，即栈顶的数
	b = self.stack.pop()

	//除了这2个为单目运算外，其余的都是双目运算，需要取第二个栈顶的数
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = self.stack.pop()
	} else {
		a = b
	}

	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		//将运算结果压入栈顶
		self.stack.push(result)
	} else {
		panic("arithmetic error! ")
	}
}

func _arith(a, b luaValue, op operator) luaValue {
	//位运算，特点是都没有浮点运算
	if op.floatFunc == nil {
		//先转成整数
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else {
		//算术运算
		if op.integerFunc != nil {
			//只有这几种算术有整数运算：+,-,*,%,//,-（取反）
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		if x, ok := convertToFloat(a); ok {
			//所有算术运算都有浮点运算
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}
	return nil
}
