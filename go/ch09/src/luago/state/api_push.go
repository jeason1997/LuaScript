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

func (self *luaState) PushNil()             { self.stack.push(nil) }
func (self *luaState) PushBoolean(b bool)   { self.stack.push(b) }
func (self *luaState) PushInteger(n int64)  { self.stack.push(n) }
func (self *luaState) PushNumber(n float64) { self.stack.push(n) }
func (self *luaState) PushString(s string)  { self.stack.push(s) }
