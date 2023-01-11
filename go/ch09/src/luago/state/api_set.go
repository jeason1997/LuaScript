/*
	 *该脚本是luago/api/lua_state.go里的接口的具体实现
	 *主要实现：Table修改方法 (stack -> Lua)
	 	SetTable(idx int)
		SetField(idx int, k string)
		SetI(idx int, n int64)
*/
package state

import . "luago/api"

//作用是把键值对写入表。其中键和值从栈里弹出，表则位于指定索引处
func (self *luaState) SetTable(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v)
}

//作用是把键值对写入表。其中键由参数传入（字符串），值从栈里弹出，表则位于指定索引处
func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, k, v)
}

//作用是把键值对写入表。其中键由参数传入（整数），值从栈里弹出，表则位于指定索引处
//用于按索引修改数组元素
func (self *luaState) SetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v)
}

func (self *luaState) setTable(t, k, v luaValue) {
	if tbl, ok := t.(*luaTable); ok {
		tbl.put(k, v)
		return
	}
	panic("not a table! ")
}

//往全局环境里写入一个值，其中字段名由参数指定，值从栈顶弹出
func (self *luaState) SetGlobal(name string) {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	v := self.stack.pop()
	self.setTable(t, name, v)
}

//专门用于给全局环境注册Go函数值。该方法仅操作全局环境，字段名和Go函数从参数传入，不改变Lua栈的状态
func (self *luaState) Register(name string, f GoFunction) {
	//先将Go函数闭包压入到栈顶
	self.PushGoFunction(f)
	//将栈顶的Go函数注册到注册表里的全局环境表里
	self.SetGlobal(name)
}
