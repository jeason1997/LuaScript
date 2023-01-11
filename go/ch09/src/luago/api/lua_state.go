package api

/*
 *我们约定，Go函数必须满足这样的签名：接收一个LuaState接口类型的参数，返回一个整数。
 *在Go函数开始执行之前，Lua栈里是传入的参数值，别无它值。
 *当Go函数结束之后，把需要返回的值留在栈顶，然后返回一个整数表示返回值个数。
 *由于Go函数返回了返回值数量，这样它在执行完毕时就不用对栈进行清理了，把返回值留在栈顶即可
 */
type GoFunction func(LuaState) int //Go函数类型
type LuaType = int                 //数据类型
type ArithOp = int                 //运算类型
type CompareOp = int               //比较类型

type LuaState interface {
	/* api_stack.go：基础栈操作方法 */

	GetTop() int             //栈顶索引，Lua从1开始
	AbsIndex(idx int) int    //把索引转换成绝对索引（并没有考虑索引是否有效）
	CheckStack(n int) bool   //检查栈的剩余空间是否能够容纳至少n个值，如果不满足，就扩容
	Pop(n int)               //从栈顶弹出一个值
	Copy(fromIdx, toIdx int) //将fromIdx处的值拷贝到toIdx处
	PushValue(idx int)       //把指定索引处的值推入栈顶
	Replace(idx int)         //将栈顶值弹出，然后写入指定位置
	Insert(idx int)          //将栈顶值弹出，然后插入指定位置
	Remove(idx int)          //删除指定索引处的值，然后将该值上面的值全部下移一个位置
	Rotate(idx, n int)       //将[idx, top]索引区间内的值朝栈顶方向旋转n个位置，如果n是负数，那么实际效果就是朝栈底方向旋转
	SetTop(idx int)          //将栈顶索引设置为指定值。如果指定值小于当前栈顶索引，效果则相当于弹出操作（指定值为0相当于清空栈）,如果指定值大于当前栈顶索引，则效果相当于推入多个nil值

	/* api_access.go：栈访问方法 (stack -> Go) */

	TypeName(tp LuaType) string        //根据类型返回该类型的字符串名称
	Type(idx int) LuaType              //根据索引返回值的类型，如果索引无效，则返回LUA_TNONE
	IsNone(idx int) bool               //是否None类型
	IsNil(idx int) bool                //是否为Nil类型
	IsNoneOrNil(idx int) bool          //是否为None或者Nil类型
	IsBoolean(idx int) bool            //是否为Boolean类型
	IsInteger(idx int) bool            //判断给定索引处的值是否是整数类型
	IsNumber(idx int) bool             //判断给定索引处的值是否是（或者可以转换为）数字类型
	IsString(idx int) bool             //判断给定索引处的值是否是字符串（或是数字）
	ToBoolean(idx int) bool            //从指定索引处取出一个布尔值，如果值不是布尔类型，则需要进行类型转换
	ToInteger(idx int) int64           //将值转换为整数类型，如果值不是整数类型并且也没办法转换成整数类型，则返回0
	ToIntegerX(idx int) (int64, bool)  //将值转换为整数类型，如果值不是整数类型并且也没办法转换成整数类型，则返回0
	ToNumber(idx int) float64          //将值转换为数字类型，如果值不是数字类型并且也没办法转换成数字类型，则返回0
	ToNumberX(idx int) (float64, bool) //将值转换为数字类型，如果值不是数字类型并且也没办法转换成数字类型，则返回0
	ToString(idx int) string           //将数值转为字符串
	ToStringX(idx int) (string, bool)  //将数值转为字符串

	/* api_push.go：压栈方法 (Go -> stack) */

	PushNil()             //将Nil值压入栈顶
	PushBoolean(b bool)   //将Boolean值压入栈顶
	PushInteger(n int64)  //将整数压入栈顶
	PushNumber(n float64) //将数字压入栈顶
	PushString(s string)  //将字符串压入栈顶

	/* api_arith.go & api_compare & api_misc.go：运算操作 */

	Arith(op ArithOp)                          //从栈顶取出操作数，按照一定规则运算，并将结果压回栈顶
	Compare(idx1, idx2 int, op CompareOp) bool //对指定索引处的两个值进行比较，返回结果。该方法不改变栈的状态
	Len(idx int)                               //访问指定索引处的值，取其长度，然后推入栈顶
	Concat(n int)                              //从栈顶弹出n个值，对这些值进行拼接，然后把结果推入栈顶

	/* api_get.go：Table访问方法 (Lua -> stack) */

	NewTable()                          //创建一个空表
	CreateTable(nArr, nRec int)         //创建一个表并推入栈顶
	GetTable(idx int) LuaType           //根据键（从栈顶弹出）从表（索引由参数指定）里取值，然后把值推入栈顶并返回值的类型
	GetField(idx int, k string) LuaType //根据键（字符串参数）从表（索引由参数指定）里取值，然后把值推入栈顶并返回值的类型
	GetI(idx int, i int64) LuaType      //根据键（数字参数）从表（索引由参数指定）里取值，然后把值推入栈顶并返回值的类型，该方法是专门给数组准备的

	/* api_set.go：Table修改方法 (stack -> Lua) */

	SetTable(idx int)           //作用是把键值对写入表。其中键和值从栈里弹出，表则位于指定索引处
	SetField(idx int, k string) //作用是把键值对写入表。其中键由参数传入（字符串），值从栈里弹出，表则位于指定索引处
	SetI(idx int, n int64)      //作用是把键值对写入表。其中键由参数传入（整数），值从栈里弹出，表则位于指定索引处，用于按索引修改数组元素

	/* api_call.go：LUA的加载与闭包的运行 */

	Load(chunk []byte, chunkName, mode string) int //从资源加载主函数原型并压入栈顶（只有主函数需要从资源加载，子函数都包括在主函数里面了）
	Call(nArgs, nResults int)                      //对Lua函数进行调用。在执行Call方法之前，必须先把被调函数推入栈顶，然后把参数值依次推入栈顶。方法结束之后，参数值和函数会被弹出栈顶，取而代之的是指定数量的返回值压入栈顶。

	/* api_push.go & api_access：Go的转换与返回 */

	PushGoFunction(f GoFunction)     //接收一个Go函数参数，把它转变成Go闭包后推入栈顶
	IsGoFunction(idx int) bool       //判断指定索引处的值是否可以转换为Go函数
	ToGoFunction(idx int) GoFunction //把指定索引处的值转换为Go函数并返回，如果值无法转换为Go函数，返回nil
}
