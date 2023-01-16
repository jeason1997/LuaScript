/*
	 *该脚本是luago/api/lua_state.go里的接口的具体实现
	 *主要实现：Table修改方法 (stack -> Lua)
	 	SetTable(idx int)
		SetField(idx int, k string)
		SetI(idx int, n int64)
*/
package state

import . "luago/api"

//作用是把键值对写入表。其中键和值从栈里弹出，表则位于指定索引处
func (self *luaState) SetTable(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v, false)
}

//SetTable的忽略元方法版本
func (self *luaState) RawSet(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v, true)
}

//作用是把键值对写入表。其中键由参数传入（字符串），值从栈里弹出，表则位于指定索引处
func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, k, v, false)
}

//作用是把键值对写入表。其中键由参数传入（整数），值从栈里弹出，表则位于指定索引处
//用于按索引修改数组元素
func (self *luaState) SetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, false)
}

//SetI的忽略元方法版本
func (self *luaState) RawSetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, true)
}

/*
 *当Lua执行t[k]=v语句时，如果t不是表，或者k在表中不存在，就会触发__newindex元方法。
 *和__index元方法一样，__newindex元方法也可以是函数或者表。
 *如果是函数，那么Lua会以t、k和v为参数调用该函数；
 *如果是表，Lua会以k为键v为值给该表赋值（可能会继续触发__newindex元方法）。
 */
func (self *luaState) setTable(t, k, v luaValue, raw bool) {
	//如果t是表
	if tbl, ok := t.(*luaTable); ok {
		//增加了raw参数，如果该参数值为true，表示需要忽略元方法。
		//如果t是表，并且键已经在表里了，或者需要忽略元方法，或者表没有__newindex元方法，则维持原来的逻辑
		if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
			tbl.put(k, v)
			return
		}
	}

	//如果raw为false，则判断有没有__newindex
	//如果t是表并且表中找不到k对应的值，而且存在__newindex元方法
	//或者t不是表
	if !raw {
		if mf := getMetafield(t, "__newindex", self); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				//如果是表，Lua会以k为键v为值给该表赋值（可能会继续触发__newindex元方法）
				self.setTable(x, k, v, false)
				return
			case *closure:
				//如果是函数，那么Lua会以t、k和v为参数调用该函数
				self.stack.push(mf)
				self.stack.push(t)
				self.stack.push(k)
				self.stack.push(v)
				self.Call(3, 0)
				return
			}
		}
	}

	panic("index error! ")
}

//往全局环境里写入一个值，其中字段名由参数指定，值从栈顶弹出
func (self *luaState) SetGlobal(name string) {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	v := self.stack.pop()
	self.setTable(t, name, v, false)
}

//专门用于给全局环境注册Go函数值。该方法仅操作全局环境，字段名和Go函数从参数传入，不改变Lua栈的状态
func (self *luaState) Register(name string, f GoFunction) {
	//先将Go函数闭包压入到栈顶
	self.PushGoFunction(f)
	//将栈顶的Go函数注册到注册表里的全局环境表里
	self.SetGlobal(name)
}

//从栈顶弹出一个表，然后把指定索引处值的元表设置成该表
func (self *luaState) SetMetatable(idx int) {
	val := self.stack.get(idx)
	mtVal := self.stack.pop()

	if mtVal == nil {
		//如果它是nil，实际效果就是清除元表
		setMetatable(val, nil, self)
	} else if mt, ok := mtVal.(*luaTable); ok {
		//如果它是表，用它设置元表
		setMetatable(val, mt, self)
	} else {
		panic("table expected! ") // todo
	}
}
