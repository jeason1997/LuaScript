package state

import . "luago/api"

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
 *Lua栈（调用帧），在执行Lua函数时，Lua栈充当虚拟寄存器以供指令操作。
 *在调用Lua/Go函数时，Lua栈充当栈帧以供参数和返回值传递。
 *一个运行于虚拟机的正常栈，里面的结构应该是，底部是程序预留的寄存器（编译时自动计算最多需要多少个寄存器Prototype.MaxStackSize）。
 *上面剩余的是计算用到的栈空间，一般会预留几个。然后初始时栈顶索引是位于寄存器上面的，也就是计算栈的初始位置。
 *slots = [reg1][reg2][reg3][stack1][stack2][...]
 *top = 4
 */
type luaStack struct {
	/* 虚拟栈 */
	slots []luaValue //栈
	top   int        //栈顶索引，Lua从1开始
	/* 函数调用信息 */
	state   *luaState        //解释器的引用
	closure *closure         //闭包（函数原型）
	varargs []luaValue       //变长参数列表
	pc      int              //程序指令地址
	openuvs map[int]*upvalue //key是寄存器索引，值是Upvalue指针
	/* 调用栈链接列表 */
	prev *luaStack //调用帧的上一个调用帧
}

//创建指定容量的栈
func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		state: state,
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
	//如果索引小于等于LUA_REGISTRYINDEX，说明是伪索引，直接返回即可
	if idx <= LUA_REGISTRYINDEX {
		return idx
	}
	if idx >= 0 {
		return idx
	}
	return idx + self.top + 1
}

//判断索引是否有效
func (self *luaStack) isValid(idx int) bool {
	//注册表伪索引属于有效索引，所以直接返回true
	if idx == LUA_REGISTRYINDEX {
		return true
	}

	//如果索引小于注册表索引，说明是Upvalue伪索引
	if idx < LUA_REGISTRYINDEX {
		//把它转成真实索引（从0开始）然后看它是否在有效范围之内
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := self.closure
		return c != nil && uvIdx < len(c.upvals)
	}

	absIdx := self.absIndex(idx)
	return absIdx > 0 && absIdx <= self.top
}

//根据索引从栈里取值，如果索引无效则返回nil值
func (self *luaStack) get(idx int) luaValue {
	//如果索引是注册表伪索引，直接返回注册表
	if idx == LUA_REGISTRYINDEX {
		return self.state.registry
	}

	//如果索引小于注册表索引，说明是Upvalue伪索引
	if idx < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := self.closure
		//如果伪索引无效，直接返回nil，否则返回Upvalue值
		if c == nil || uvIdx >= len(c.upvals) {
			return nil
		}
		return *(c.upvals[uvIdx].val)
	}

	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		return self.slots[absIdx-1]
	}
	return nil
}

//根据索引往栈里写入值，如果索引无效，则调用panic()函数终止程序
func (self *luaStack) set(idx int, val luaValue) {
	//如果索引是注册表伪索引，直接修改注册表
	if idx == LUA_REGISTRYINDEX {
		self.state.registry = val.(*luaTable)
	}

	//如果索引小于注册表索引，说明是Upvalue伪索引
	if idx < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := self.closure
		//如果伪索引有效，我们就修改Upvalue值，否则直接返回
		if c != nil && uvIdx < len(c.upvals) {
			*(c.upvals[uvIdx].val) = val
		}
	}

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
