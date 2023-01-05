package vm

/*
 *NEWTABLE指令是iABC模式，操作数B和C只有9个比特，如果当作无符号整数的话，最大也不能超过512。
 *但是我们在前面也提到过，因为表构造器便捷实用，所以Lua也经常被用来描述数据（类似JSON），
 *如果有很大的数据需要写成表构造器，但是表的初始容量又不够大，就容易导致表频繁扩容从而影响数据加载效率。
 *为了解决这个问题，NEWTABLE指令的B和C操作数使用了一种叫作浮点字节（Floating Point Byte）的编码方式。
 *这种编码方式和浮点数的编码方式类似，只是仅用一个字节。
 *具体来说，如果把某个字节用二进制写成eeeeexxx，那么当eeeee == 0时该字节表示的整数就是xxx，
 *否则该字节表示的整数是(1xxx) ＊ 2^(eeeee -1)。
 */

//将整数转化为“浮点字节”编码方式
func Int2fb(x int) int {
	//指数
	e := 0

	//00000111
	if x < 8 {
		return x
	}
	for x >= (8 << 4) { /* coarse steps */
		x = (x + 0xf) >> 4 /* x = ceil(x / 16) */
		e += 4
	}
	for x >= (8 << 1) { /* fine steps */
		x = (x + 1) >> 1 /* x = ceil(x / 2) */
		e++
	}
	return ((e + 1) << 3) | (x - 8)
}

//将“浮点字节”编码方式的数转为整数
func Fb2int(x int) int {
	if x < 8 {
		return x
	} else {
		return ((x & 7) + 8) << uint((x>>3)-1)
	}
}
