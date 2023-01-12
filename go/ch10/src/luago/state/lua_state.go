package state

import . "luago/api"

//Lua解释器
type luaState struct {
	/*
	 *Lua给用户提供了一个注册表，这个注册表实际上就是一个普通的Lua表，所以用户可以在里面存放任何Lua值。
	 *有趣的是，这个注册表虽然是给用户准备的，但Lua本身也用到了它，比如说Lua全局变量就是借助这个注册表实现的。
	 *由于注册表是全局状态，每个Lua解释器实例都有自己的注册表，所以把它放在luaState结构体里是合理的。
	 *访问：由于注册表实际就是个普通的Lua表，所以Lua API并没有提供专门的方法来操作注册表。
	 *任何可以操作表的API方法（比如GetTable()等）都可以用来操作注册表。
	 *但是，普遍表都是存在luaStack里的，操作方法都是通过索引来访问表的，而注册表是特殊存在luaState里的，那么怎么访问注册表呢？
	 *答案是通过“伪索引（pseudo-index）”，只要访问表的时候传入的索引是LUA_REGISTRYINDEX，就会自动去访问呢注册表。
	 */
	registry *luaTable

	/*
	 *Lua栈（调用帧），在执行Lua函数时，Lua栈充当虚拟寄存器以供指令操作。
	 *在调用Lua/Go函数时，Lua栈充当栈帧以供参数和返回值传递。
	 */
	stack *luaStack
}

func New() *luaState {
	//先创建注册表
	registry := newLuaTable(0, 0)
	//然后预先往里面放一个全局环境，所有的Lua全局变量都放在这个表里
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0))

	ls := &luaState{registry: registry}
	//推入一个空的Lua栈（调用帧）
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
	return ls
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

//对于任何一个Upvalue索引，用注册表伪索引减去该索引就可以得到对应的Upvalue伪索引
//在Lua虚拟机指令的操作数里，Upvalue索引是从0开始的，但是在转换成Lua栈伪索引时，Upvalue指令是从1开始的
func LuaUpvalueIndex(i int) int {
	return LUA_REGISTRYINDEX - i
}
