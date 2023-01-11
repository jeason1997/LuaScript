package state

type luaState struct {
	stack *luaStack
}

func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
}

// 压入一个调用帧
func (self *luaState) pushLuaStack(stack *luaStack) {
	//使用单向链表的方式实现函数调用栈
	//往栈顶推入一个调用帧相当于在链表头部插入一个节点，并让这个节点成为新的头部
	stack.prev = self.stack
	self.stack = stack
}

// 从栈顶弹出一个调用帧
func (self *luaState) popLuaStack() {
	stack := self.stack
	//将栈顶帧改为链接的上一个调用帧
	self.stack = stack.prev
	//原栈顶帧断开连接
	stack.prev = nil
}
