package state

import (
	"fmt"
	. "luago/api"
	"luago/number"
)

type luaValue interface{}

func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64:
		return LUA_TNUMBER
	case float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	case *luaTable:
		return LUA_TTABLE
	case *closure:
		return LUA_TFUNCTION
	default:
		panic("todo! ")
	}
}

// 在Lua里，只有false和nil表示假，其他一切值都表示真
func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

// 任意数值转浮点数
func convertToFloat(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		return number.ParseFloat(x)
	default:
		return 0, false
	}
}

// 任意数值转整数
func convertToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		return _stringToInteger(x)
	default:
		return 0, false
	}
}

func _stringToInteger(s string) (int64, bool) {
	if i, ok := number.ParseInteger(s); ok {
		return i, true
	}
	if f, ok := number.ParseFloat(s); ok {
		return number.FloatToInteger(f)
	}
	return 0, false
}

// 设置元表，如果是Table的话，直接设置，否则到注册表里设置共享元表
func setMetatable(val luaValue, mt *luaTable, ls *luaState) {
	//如果是表的话，直接设置元表
	if t, ok := val.(*luaTable); ok {
		t.metatable = mt
		return
	}
	//虽然注册表也是一个普通的表，不过按照约定，下划线开头后跟大写字母的字段名是保留给Lua实现使用的，
	//所以我们使用了“_MT1”这样的字段名，以免和用户（通过API）放在注册表里的数据产生冲突
	key := fmt.Sprintf("_MT%d", typeOf(val))
	//其他值则是每种类型共享一个元表，放在注册表里
	ls.registry.put(key, mt)
}

// 获取元表，如果是Table的话，直接获取，否则到注册表里获取共享元表
func getMetatable(val luaValue, ls *luaState) *luaTable {
	if t, ok := val.(*luaTable); ok {
		return t.metatable
	}
	key := fmt.Sprintf("_MT%d", typeOf(val))
	if mt := ls.registry.get(key); mt != nil {
		return mt.(*luaTable)
	}
	return nil
}

// 获取某种类型的变量的某个元方法
func getMetafield(val luaValue, fieldName string, ls *luaState) luaValue {
	if mt := getMetatable(val, ls); mt != nil {
		return mt.get(fieldName)
	}
	return nil
}

func callMetamethod(a, b luaValue, mmName string, ls *luaState) (luaValue, bool) {
	var mm luaValue
	//两个操作数，只要其中一个有元方法就行
	if mm = getMetafield(a, mmName, ls); mm == nil {
		if mm = getMetafield(b, mmName, ls); mm == nil {
			return nil, false
		}
	}

	//如果任何一个操作数有对应元方法，则以两个操作数为参数调用元方法，将元方法调用结果和true返回
	ls.stack.check(4)
	//压入元方法
	ls.stack.push(mm)
	//压入两个操作数
	ls.stack.push(a)
	ls.stack.push(b)
	//调用元方法，2个参数，一个返回值
	ls.Call(2, 1)
	//将栈顶的返回值弹出，恢复栈状态
	return ls.stack.pop(), true
}
