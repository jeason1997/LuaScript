/*
 *该脚本是luago/api/lua_state.go里的接口的具体实现
 *主要实现：运算操作
 *	Compare(idx1, idx2 int, op CompareOp) bool
 */
package state

import (
	. "luago/api"
)

// 与Compare的LUA_OPEQ逻辑大体相同
// 但是当值为Table时，只进行基本的比较，不调用Table的__eq元方法比较
func (self *luaState) RawEqual(idx1, idx2 int) bool {
	if !self.stack.isValid(idx1) || !self.stack.isValid(idx2) {
		return false
	}

	a := self.stack.get(idx1)
	b := self.stack.get(idx2)
	return _eq(a, b, nil)
}

// 对指定索引处的两个值进行比较，返回结果。该方法不改变栈的状态
func (self *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a := self.stack.get(idx1)
	b := self.stack.get(idx2)
	switch op {
	case LUA_OPEQ:
		return _eq(a, b, self)
	case LUA_OPLT:
		return _lt(a, b, self)
	case LUA_OPLE:
		return _le(a, b, self)
	default:
		panic("invalid compare op! ")
	}
}

// 等于
func _eq(a, b luaValue, ls *luaState) bool {
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
	case *luaTable:
		//当且仅当两个操作数是不同的表时，才会尝试执行__eq元方法
		if y, ok := b.(*luaTable); ok && x != y && ls != nil {
			if result, ok := callMetamethod(x, y, "__eq", ls); ok {
				//执行结果会被转换为布尔值
				return convertToBoolean(result)
			}
		}
		return a == b
	default:
		return a == b
	}
}

// 小于
func _lt(a, b luaValue, ls *luaState) bool {
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

	//如果不是数字也不是字符串，则尝试元方法
	if result, ok := callMetamethod(a, b, "__lt", ls); ok {
		return convertToBoolean(result)
	} else {
		panic("comparison error! ")
	}
}

// 小于等于
func _le(a, b luaValue, ls *luaState) bool {
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

	//如果不是数字也不是字符串，则尝试元方法
	if result, ok := callMetamethod(a, b, "__le", ls); ok {
		return convertToBoolean(result)
	} else if result, ok := callMetamethod(b, a, "__lt", ls); ok {
		//如果Lua找不到__le元方法，则会尝试调用__lt元方法（假设a <= b等价于not (b < a)）
		return !convertToBoolean(result)
	} else {
		panic("comparison error! ")
	}
}
