/*
 *该脚本是luago/api/lua_state.go里的接口的具体实现
 *主要实现：运算操作
 *	Compare(idx1, idx2 int, op CompareOp) bool
 */
package state

import . "luago/api"

//对指定索引处的两个值进行比较，返回结果。该方法不改变栈的状态
func (self *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a := self.stack.get(idx1)
	b := self.stack.get(idx2)
	switch op {
	case LUA_OPEQ:
		return _eq(a, b)
	case LUA_OPLT:
		return _lt(a, b)
	case LUA_OPLE:
		return _le(a, b)
	default:
		panic("invalid compare op! ")
	}
}

//等于
func _eq(a, b luaValue) bool {
	//只有当两个操作数在Lua语言层面具有相同类型时，等于运算才有可能返回true。
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		//整数和浮点数仅仅在Lua实现层面有差别，在Lua语言层面统一表现为数字类型，因此需要相互转换
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	default:
		return a == b
	}
}

//小于
func _lt(a, b luaValue) bool {
	//小于操作仅对数字和字符串类型有意义
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		}
	}
	panic("comparison error! ")
}

//小于等于
func _le(a, b luaValue) bool {
	//小于等于操作仅对数字和字符串类型有意义
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}
	panic("comparison error! ")
}
