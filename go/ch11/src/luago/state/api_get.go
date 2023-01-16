/*
	 *该脚本是luago/api/lua_state.go里的接口的具体实现
	 *主要实现：Table访问方法 (Lua -> stack)
	 	NewTable()
		CreateTable(nArr, nRec int)
		GetTable(idx int) LuaType
		GetField(idx int, k string) LuaType
		GetI(idx int, i int64) LuaType
*/
package state

import . "luago/api"

//创建一个表并推入栈顶
func (self *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	self.stack.push(t)
}

//创建一个空表
func (self *luaState) NewTable() {
	self.CreateTable(0, 0)
}

//根据键（从栈顶弹出）从表（索引由参数指定）里取值，然后把值推入栈顶并返回值的类型
func (self *luaState) GetTable(idx int) LuaType {
	t := self.stack.get(idx)
	k := self.stack.pop()
	return self.getTable(t, k, false)
}

//GetTable的忽略元方法版本
func (self *luaState) RawGet(idx int) LuaType {
	t := self.stack.get(idx)
	k := self.stack.pop()
	return self.getTable(t, k, true)
}

//根据键（字符串参数）从表（索引由参数指定）里取值，然后把值推入栈顶并返回值的类型
//和GetTable()方法类似，只不过键不是从栈顶弹出的任意值，而是由参数传入的字符串
func (self *luaState) GetField(idx int, k string) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, k, false)
}

//根据键（数字参数）从表（索引由参数指定）里取值，然后把值推入栈顶并返回值的类型
//和GetField()方法类似，只不过由参数传入的键是数字而非字符串，该方法是专门给数组准备的
func (self *luaState) GetI(idx int, i int64) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, i, false)
}

//GeI的忽略元方法版本
func (self *luaState) RawGetI(idx int, i int64) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, i, true)
}

/*
 *当Lua执行t[k]表达式时，如果t不是表，或者k在表中不存在，
 *就会触发__index元方法。虽然名为元方法，但实际上__index元方法既可以是函数，也可以是表。
 *如果是函数，那么Lua会以t和k为参数调用该函数，以函数返回值为结果；
 *如果是表，Lua会以k为键访问该表，以值为结果（可能会继续触发__index元方法）
 */
func (self *luaState) getTable(t, k luaValue, raw bool) LuaType {
	//如果t是表
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		//增加了raw参数，如果该参数值为true，表示需要忽略元方法。
		//如果t是表，并且键已经在表里了，或者需要忽略元方法，或者表没有__index元方法，则维持原来的逻辑
		if raw || v != nil || !tbl.hasMetafield("__index") {
			self.stack.push(v)
			return typeOf(v)
		}
	}

	//如果raw为false，则判断有没有__index
	//如果t是表并且表中找不到k对应的值，而且存在__index元方法
	//或者t不是表
	if !raw {
		if mf := getMetafield(t, "__index", self); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				//如果元方法是一个Table，Lua会以k为键访问该表，以值为结果（可能会继续触发__index元方法）
				return self.getTable(x, k, false)
			case *closure:
				//如果是函数，则以t和k为参数调用该函数
				self.stack.push(mf)
				self.stack.push(t)
				self.stack.push(k)
				self.Call(2, 1)
				v := self.stack.get(-1)
				return typeOf(v)
			}
		}
	}

	panic("index error!")
}

//把全局环境中的某个字段（名字由参数指定）推入栈顶
func (self *luaState) GetGlobal(name string) LuaType {
	t := self.registry.get(LUA_RIDX_GLOBALS)
	return self.getTable(t, name, false)
}

//看指定索引处的值是否有元表，如果有，则把元表推入栈顶并返回true；否则栈的状态不改变，返回false。
func (self *luaState) GetMetatable(idx int) bool {
	val := self.stack.get(idx)

	if mt := getMetatable(val, self); mt != nil {
		self.stack.push(mt)
		return true
	} else {
		return false
	}
}
