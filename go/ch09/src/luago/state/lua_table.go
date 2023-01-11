/*
 *表：关联数组，里面存放的是两两关联的键值对。除了nil值和浮点数NaN以外，任何Lua值都可以当作键来使用。值则可以是任意Lua值，包括nil和NaN。
 *记录：表的键全部是字符串
 *数组：表的键全部是正整数
 *序列：数组中不存在nil值
 *
 *表
 *  -记录
 *  -数组
 *	  -序列
 */
package state

import (
	"luago/number"
	"math"
)

/*
 *参考Lua官方的做法，使用数组和哈希表的混合方式来实现Lua表。
 *如果表的键是连续的正整数，那么哈希表就是空的，值全部按索引存储在数组里。
 */
type luaTable struct {
	arr []luaValue
	//由于map是Go语言关键字，不能用来命名字段，所以加了下划线
	_map map[luaValue]luaValue
}

// 该函数接受两个参数，用于预估表的用途和容量。
// 如果参数nArr大于0，说明表可能是当作数组使用的，先创建数组部分；如果参数nRec大于0，说明表可能是当作记录使用的，先创建哈希表部分
func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}
	return t
}

// 根据键从表里查找值。
func (self *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		//如果键是整数（或者能够转换为整数的浮点数），且在数组索引范围之内，直接按索引访问数组部分就可以了；
		if idx >= 1 && idx <= int64(len(self.arr)) {
			return self.arr[idx-1]
		}
	}
	// 否则从哈希表查找值。
	return self._map[key]
}

// 往表里存入键值对
func (self *luaTable) put(key, val luaValue) {
	//不允许用nil或者NaN作为key
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}

	key = _floatToInteger(key)
	//值如果是整数，并且值在Array里，或者Array的尾部+1，则放到Array里，否则当成Map数据
	if idx, ok := key.(int64); ok && idx >= 1 {
		//如果键是（或者已经被转换为）整数，且在数组索引范围之内的话，直接按索引修改数组元素就可以了
		arrLen := int64(len(self.arr))
		if idx <= arrLen {
			self.arr[idx-1] = val
			//向数组里放入nil值会制造洞，如果洞在数组末尾的话，调用_shrinkArray()函数把尾部的洞全部删除
			if idx == arrLen && val == nil {
				self._shrinkArray()
			}
			return
		}
		//如果键是整数，而且刚刚超出数组索引范围且值不是nil，就把值追加到数组末尾，然后调用_expandArray()函数动态扩展数组
		if idx == arrLen+1 {
			//把原本存在哈希表里的某些值也挪到数组里
			delete(self._map, key)
			if val != nil {
				//扩容追加到数组末尾
				self.arr = append(self.arr, val)
				//重新调整Array，把之前放在Map里的，条件满足（key的值刚好是在Array尾部+1）的挪到Array里
				self._expandArray()
			}
			return
		}
	}
	//如果值不是nil就把键值对写入哈希表，否则把键从哈希表里删除以节约空间
	if val != nil {
		if self._map == nil {
			//由于在创建表的时候并不一定创建了哈希表部分，所以在第一次写入时，需要创建哈希表
			self._map = make(map[luaValue]luaValue, 8)
		}
		self._map[key] = val
	} else {
		delete(self._map, key)
	}
}

// Table的len只对其Array部分有效，不能用于计算Map长度
func (self *luaTable) len() int {
	return len(self.arr)
}

// 尝试把浮点数类型的键转换成整数
func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

// 删除数组里尾部的洞
func (self *luaTable) _shrinkArray() {
	for i := len(self.arr) - 1; i >= 0; i-- {
		if self.arr[i] == nil {
			self.arr = self.arr[0:i]
		} else {
			break
		}
	}
}

/*
 *动态扩展数组，把存放在Map里的，且key是正整数，并且Key数值是在Array的当前容量之后的连续数据，从Map转为Array
 *比如当前的Array长度是5，Map里有{6:value, 7:value, 9:value}，则把6，7这两个数据挪到Array里
 *为什么这样做呢？因为往table里存放数据的时候，如果Key是正整数，那么会判断Key是否在Array里或Array的长度+1，
 *是的话就放进Array，否则当成Map数据。随着Array的扩展，之前放在Map里的零散数据，key值刚好可以转到Array里，就动态挪过去
 */
func (self *luaTable) _expandArray() {
	for idx := int64(len(self.arr)) + 1; true; idx++ {
		if val, found := self._map[idx]; found {
			delete(self._map, idx)
			self.arr = append(self.arr, val)
		} else {
			break
		}
	}
}
