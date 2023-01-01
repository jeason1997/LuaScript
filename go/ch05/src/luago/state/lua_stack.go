package state

/*
 *栈容量：n
 *栈顶索引：top(0 < top <= n)
 *绝对索引：正数索引叫作绝对索引，从1（栈底）开始递增
 *相对索引：负数索引叫作相对索引，从-1（栈顶）开始递减
 *有效索引：位于[1, top]区间的叫有效索引
 *无效索引：位于(top, n]区间的叫无效索引
 *可接受索引：位于[1, n]区间的叫可接受索引
 */

type luaStack struct {
	slots []luaValue
	top   int //栈顶索引，Lua从1开始
}

//创建指定容量的栈
func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}

//检查栈的剩余空间是否能够容纳至少n个值，如果不满足，就扩容
func (self *luaStack) check(n int) {
	free := len(self.slots) - self.top
	for i := free; i < n; i++ {
		self.slots = append(self.slots, nil)
	}
}

//将值推入栈
func (self *luaStack) push(val luaValue) {
	if self.top == len(self.slots) {
		panic("stack overflow! ")
	}
	self.slots[self.top] = val
	self.top++
}

//从栈顶弹出一个值
func (self *luaStack) pop() luaValue {
	if self.top < 1 {
		panic("stack underflow! ")
	}
	self.top--
	val := self.slots[self.top]
	self.slots[self.top] = nil
	return val
}

//把索引转换成绝对索引（并没有考虑索引是否有效）
func (self *luaStack) absIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	return idx + self.top + 1
}

//判断索引是否有效
func (self *luaStack) isValid(idx int) bool {
	absIdx := self.absIndex(idx)
	return absIdx > 0 && absIdx <= self.top
}

//根据索引从栈里取值，如果索引无效则返回nil值
func (self *luaStack) get(idx int) luaValue {
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		return self.slots[absIdx-1]
	}
	return nil
}

//根据索引往栈里写入值，如果索引无效，则调用panic()函数终止程序
func (self *luaStack) set(idx int, val luaValue) {
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		self.slots[absIdx-1] = val
		return
	}
	panic("invalid index! ")
}

//将from到to索引范围内的数据反转
func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
