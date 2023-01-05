/*
	 *该脚本是luago/api/lua_state.go里的接口的具体实现
	 *主要实现：Table修改方法 (stack -> Lua)
	 	SetTable(idx int)
		SetField(idx int, k string)
		SetI(idx int, n int64)
*/
package state

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
