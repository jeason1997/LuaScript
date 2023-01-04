/*
	 *该脚本是luago/api/lua_state.go里的接口的具体实现
	 *主要实现：运算操作
	 	Len(idx int)
		Concat(n int)
*/
package state

//访问指定索引处的值，取其长度，然后推入栈顶
func (self *luaState) Len(idx int) {
	val := self.stack.get(idx)
	//暂时只考虑字符串的长度，对于其他情况则调用panic()函数终止程序
	if s, ok := val.(string); ok {
		self.stack.push(int64(len(s)))
	} else {
		panic("length error! ")
	}
}

//从栈顶弹出n个值，对这些值进行拼接，然后把结果推入栈顶
func (self *luaState) Concat(n int) {
	if n == 0 {
		self.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if self.IsString(-1) && self.IsString(-2) {
				s2 := self.ToString(-1)
				s1 := self.ToString(-2)
				self.stack.pop()
				self.stack.pop()
				self.stack.push(s1 + s2)
				continue
			}
			panic("concatenation error! ")
		}
	}
	// n == 1, do nothing
}
