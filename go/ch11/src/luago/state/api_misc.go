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

	if str, ok := val.(string); ok {
		//先判断值是否是字符串，如果是，结果就是字符串长度
		self.stack.push(int64(len(str)))
	} else if result, ok := callMetamethod(val, val, "__len", self); ok {
		//否则看值是否有__len元方法，如果有，则以值为参数调用元方法，将元方法返回值作为结果
		self.stack.push(result)
	} else if t, ok := val.(*luaTable); ok {
		//如果找不到对应元方法，但值是表，结果就是表的长度
		self.stack.push(int64(t.len()))
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
			//字符串拼接
			if self.IsString(-1) && self.IsString(-2) {
				s2 := self.ToString(-1)
				s1 := self.ToString(-2)
				self.stack.pop()
				self.stack.pop()
				self.stack.push(s1 + s2)
				continue
			}

			//不是字符串则调用元方法拼接
			b := self.stack.pop()
			a := self.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", self); ok {
				self.stack.push(result)
				continue
			}

			panic("concatenation error!")
		}
	}
	// n == 1, do nothing
}
