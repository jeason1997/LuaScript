package api

type LuaVM interface {
	/* api_vm.go：虚拟机的基础操作 */

	LuaState
	PC() int            // 返回当前PC（仅测试用）
	AddPC(n int)        // 修改PC（用于实现跳转指令）
	Fetch() uint32      // 根据PC索引从函数原型的指令表里取出当前指令，然后把PC加1，这样下次再调用该方法取出的就是下一条指令
	GetConst(idx int)   // 根据索引从函数原型的常量表里取出一个常量值，然后把它推入栈顶
	GetRK(rk int)       // 根据情况调用GetConst()方法把某个常量推入栈顶，或者调用PushValue()方法把某个索引处的栈值推入栈顶
	RegisterCount() int // 返回当前Lua函数所操作的寄存器数量（编译时自动计算最大需要数量）
	LoadVararg(n int)   // 把传递给当前Lua函数的变长参数推入栈顶（多退少补）
	LoadProto(idx int)  // 把当前Lua函数的子函数的原型实例化为闭包推入栈顶
}
