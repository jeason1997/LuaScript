/*
	 *该脚本是luago/api/lua_state.go里的接口的具体实现
	 *主要实现：基础栈操作方法
	 	GetTop() int
		AbsIndex(idx int) int
		CheckStack(n int) bool
		Pop(n int)
		Copy(fromIdx, toIdx int)
		PushValue(idx int)
		Replace(idx int)
		Insert(idx int)
		Remove(idx int)
		Rotate(idx, n int)
		SetTop(idx int)
*/
package state

func (self *luaState) GetTop() int {
	return self.stack.top
}

func (self *luaState) AbsIndex(idx int) int {
	return self.stack.absIndex(idx)
}

func (self *luaState) CheckStack(n int) bool {
	self.stack.check(n)
	return true // never fails
}

func (self *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
}

func (self *luaState) Copy(fromIdx, toIdx int) {
	val := self.stack.get(fromIdx)
	self.stack.set(toIdx, val)
}

//把指定索引处的值推入栈顶
func (self *luaState) PushValue(idx int) {
	val := self.stack.get(idx)
	self.stack.push(val)
}

//是PushValue()的反操作：将栈顶值弹出，然后写入指定位置
func (self *luaState) Replace(idx int) {
	val := self.stack.pop()
	self.stack.set(idx, val)
}

//将栈顶值弹出，然后插入指定位置
func (self *luaState) Insert(idx int) {
	//可以理解为从idx开始朝栈顶旋转一个单位
	self.Rotate(idx, 1)
}

//删除指定索引处的值，然后将该值上面的值全部下移一个位置
func (self *luaState) Remove(idx int) {
	//可以理解为从idx开始朝栈底方向旋转，然后删除最顶端的值
	self.Rotate(idx, -1)
	self.Pop(1)
}

//将[idx, top]索引区间内的值朝栈顶方向旋转n个位置，如果n是负数，那么实际效果就是朝栈底方向旋转
//所谓的旋转，可以理解为把从栈顶到idx的元素依次朝某个方向（栈顶或栈底）移动n个单位，移动后超出idx或者栈顶方向的，则循环插到栈顶或者idx处
func (self *luaState) Rotate(idx, n int) {
	t := self.stack.top - 1
	p := self.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	self.stack.reverse(p, m)
	self.stack.reverse(m+1, t)
	self.stack.reverse(p, t)
}

//将栈顶索引设置为指定值。如果指定值小于当前栈顶索引，效果则相当于弹出操作（指定值为0相当于清空栈）,如果指定值大于当前栈顶索引，则效果相当于推入多个nil值
func (self *luaState) SetTop(idx int) {
	newTop := self.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow! ")
	}

	n := self.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			self.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			self.stack.push(nil)
		}
	}
}
