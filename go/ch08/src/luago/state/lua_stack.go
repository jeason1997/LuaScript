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

/*
 *一个运行于虚拟机的正常栈，里面的结构应该是，底部是程序预留的寄存器（编译时自动计算最多需要多少个寄存器Prototype.MaxStackSize）。
 *上面剩余的是计算用到的栈空间，一般会预留几个。然后初始时栈顶索引是位于寄存器上面的，也就是计算栈的初始位置。
 *slots = [reg1][reg2][reg3][stack1][stack2][...]
 *top = 4
 */
type luaStack struct {
	/* 虚拟栈 */
	slots []luaValue
	top   int //栈顶索引，Lua从1开始
	/* 函数调用信息 */
	closure *closure   //闭包（函数原型）
	varargs []luaValue //变长参数列表
	pc      int        //程序指令地址
	/* 调用栈链接列表 */
	prev *luaStack //调用栈的上一个调用帧
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

//往栈顶推入多个值（多退少补）
func (self *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}
	for i := 0; i < n; i++ {
		if i < nVals {
			self.push(vals[i])
		} else {
			//如果n大于vals的长度，则后面都用nil补充
			self.push(nil)
		}
	}
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

//从栈顶一次性弹出多个值
func (self *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = self.pop()
	}
	return vals
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

//压入一个调用帧
func (self *luaState) pushLuaStack(stack *luaStack) {
	//使用单向链表的方式实现函数调用栈
	//往栈顶推入一个调用帧相当于在链表头部插入一个节点，并让这个节点成为新的头部
	stack.prev = self.stack
	self.stack = stack
}

//从栈顶弹出一个调用帧
func (self *luaState) popLuaStack() {
	stack := self.stack
	//将栈顶帧改为链接的上一个调用帧
	self.stack = stack.prev
	//原栈顶帧断开连接
	stack.prev = nil
}
