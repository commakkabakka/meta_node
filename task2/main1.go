package main

import (
	"fmt"
	"unsafe"
)

/*
	题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
*/

func ModPtr(ptr **int) {
	// （不推荐了）在 Go 的内存模型中，uintptr 不是稳定的指针类型，一旦你将指针转为整数类型（uintptr），Go GC 不再追踪这个指针的对象。
	//var p uintptr = uintptr(unsafe.Pointer(*ptr))
	//p = p + 10
	//*ptr = (*int)(unsafe.Pointer(p))

	// (推荐） 官方明确说明：If you convert a pointer to uintptr, use it only for the immediate conversion back.
	// uintptr 不在中间阶段被保存到变量中。这样 GC 不会有机会在中间阶段运行。
	*ptr = (*int)(unsafe.Add(unsafe.Pointer(*ptr), 10))
}

func Test21() {
	var a int = 10
	var p *int = &a
	fmt.Println("修改前的地址：", p)
	ModPtr(&p)
	fmt.Println("修改后的地址：", p)
}
