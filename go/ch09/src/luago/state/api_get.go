/*
	 *该脚本是luago/api/lua_state.go里的接口的具体实现
	 *主要实现：Table访问方法 (Lua -> stack)
	 	NewTable()
		CreateTable(nArr, nRec int)
		GetTable(idx int) LuaType
		GetField(idx int, k string) LuaType
		GetI(idx int, i int64) LuaType
*/
package state

import . "luago/api"

//创建一个表并推入栈顶
func (self *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	self.stack.push(t)
}

//创建一个空表
func (self *luaState) NewTable() {
	self.CreateTable(0, 0)
}

//根据键（从栈顶弹出）从表（索引由参数指定）里取值，然后把值推入栈顶并返回值的类型
func (self *luaState) GetTable(idx int) LuaType {
	t := self.stack.get(idx)
	k := self.stack.pop()
	return self.getTable(t, k)
}

//根据键（字符串参数）从表（索引由参数指定）里取值，然后把值推入栈顶并返回值的类型
//和GetTable()方法类似，只不过键不是从栈顶弹出的任意值，而是由参数传入的字符串
func (self *luaState) GetField(idx int, k string) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, k)
}

//根据键（数字参数）从表（索引由参数指定）里取值，然后把值推入栈顶并返回值的类型
//和GetField()方法类似，只不过由参数传入的键是数字而非字符串，该方法是专门给数组准备的
func (self *luaState) GetI(idx int, i int64) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, i)
}

func (self *luaState) getTable(t, k luaValue) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		self.stack.push(v)
		return typeOf(v)
	}
	panic("not a table! ")
}

//把全局环境中的某个字段（名字由参数指定）推入栈顶
func (self *luaState) GetGlobal(name string) LuaType {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	return self.getTable(t, name)
}
