/*
	 *该脚本是luago/api/lua_state.go里的接口的具体实现
	 *主要实现：栈访问方法 (stack -> Go)
	 	TypeName(tp LuaType) string
		Type(idx int) LuaType
		IsNone(idx int) bool
		IsNil(idx int) bool
		IsNoneOrNil(idx int) bool
		IsBoolean(idx int) bool
		IsInteger(idx int) bool
		IsNumber(idx int) bool
		IsString(idx int) bool
		ToBoolean(idx int) bool
		ToInteger(idx int) int64
		ToIntegerX(idx int) (int64, bool)
		ToNumber(idx int) float64
		ToNumberX(idx int) (float64, bool)
		ToString(idx int) string
		ToStringX(idx int) (string, bool)
*/
package state

import (
	"fmt"
	. "luago/api"
)

func (self *luaState) RawLen(idx int) uint {
	val := self.stack.get(idx)
	switch x := val.(type) {
	case string:
		return uint(len(x))
	case *luaTable:
		return uint(x.len())
	default:
		return 0
	}
}

func (self *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// 根据索引返回值的类型，如果索引无效，则返回LUA_TNONE
func (self *luaState) Type(idx int) LuaType {
	if self.stack.isValid(idx) {
		val := self.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

func (self *luaState) IsNone(idx int) bool {
	return self.Type(idx) == LUA_TNONE
}

func (self *luaState) IsNil(idx int) bool {
	return self.Type(idx) == LUA_TNIL
}

func (self *luaState) IsNoneOrNil(idx int) bool {
	return self.Type(idx) <= LUA_TNIL
}

func (self *luaState) IsBoolean(idx int) bool {
	return self.Type(idx) == LUA_TBOOLEAN
}

// 判断给定索引处的值是否是字符串（或是数字）
func (self *luaState) IsString(idx int) bool {
	t := self.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

// 判断给定索引处的值是否是（或者可以转换为）数字类型
func (self *luaState) IsNumber(idx int) bool {
	_, ok := self.ToNumberX(idx)
	return ok
}

// 判断给定索引处的值是否是整数类型
func (self *luaState) IsInteger(idx int) bool {
	val := self.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

// 从指定索引处取出一个布尔值，如果值不是布尔类型，则需要进行类型转换
func (self *luaState) ToBoolean(idx int) bool {
	val := self.stack.get(idx)
	return convertToBoolean(val)
}

func (self *luaState) ToNumber(idx int) float64 {
	n, _ := self.ToNumberX(idx)
	return n
}

// 将值转换为数字类型，如果值不是数字类型并且也没办法转换成数字类型，则返回0
func (self *luaState) ToNumberX(idx int) (float64, bool) {
	val := self.stack.get(idx)
	return convertToFloat(val)
}

func (self *luaState) ToInteger(idx int) int64 {
	i, _ := self.ToIntegerX(idx)
	return i
}

func (self *luaState) ToIntegerX(idx int) (int64, bool) {
	val := self.stack.get(idx)
	return convertToInteger(val)
}

func (self *luaState) ToString(idx int) string {
	s, _ := self.ToStringX(idx)
	return s
}

func (self *luaState) ToStringX(idx int) (string, bool) {
	val := self.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		//如果值是数字，则将值转换为字符串（注意会修改栈）
		s := fmt.Sprintf("%v", x)
		self.stack.set(idx, s) // 注意这里会修改栈！
		return s, true
	default:
		//其他返回空字符串
		return "", false
	}
}

// 判断指定索引处的值是否可以转换为Go函数
func (self *luaState) IsGoFunction(idx int) bool {
	val := self.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc != nil
	}
	return false
}

// 把指定索引处的值转换为Go函数并返回，如果值无法转换为Go函数，返回nil
func (self *luaState) ToGoFunction(idx int) GoFunction {
	val := self.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc
	}
	return nil
}
