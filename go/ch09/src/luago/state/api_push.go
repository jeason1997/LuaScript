/*
	 *该脚本是luago/api/lua_state.go里的接口的具体实现
	 *主要实现：压栈方法 (Go -> stack)
	 	PushNil()
		PushBoolean(b bool)
		PushInteger(n int64)
		PushNumber(n float64)
		PushString(s string)
*/
package state

import . "luago/api"

func (self *luaState) PushNil()                    { self.stack.push(nil) }
func (self *luaState) PushBoolean(b bool)          { self.stack.push(b) }
func (self *luaState) PushInteger(n int64)         { self.stack.push(n) }
func (self *luaState) PushNumber(n float64)        { self.stack.push(n) }
func (self *luaState) PushString(s string)         { self.stack.push(s) }
func (self *luaState) PushGoFunction(f GoFunction) { self.stack.push(newGoClosure(f)) }

//由于全局环境也只是个普通的Lua表，所以GetTable()和SetTable()等表操作方法也同样适用于它，
//不过要使用这些方法，必须把全局环境预先放入栈里。
//PushGlobalTable()方法就是用来做这件事的，它把全局环境推入栈顶以备后续操作使用
func (self *luaState) PushGlobalTable() {
	global := self.registry.get(LUA_RIDX_GLOBALS)
	self.stack.push(global)
}
