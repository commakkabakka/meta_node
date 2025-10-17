package main

import "fmt"

/*
	题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
*/

func ModValue(slice []int) {
	// (错误) value 是 切片中元素的副本（拷贝），不是原元素的引用。
	//for _, value := range slice {
	//	value = value * 2
	//}

	for index := range slice {
		slice[index] = slice[index] * 2
	}
}

func Test22() {
	var s1 []int = []int{1, 2, 3, 4, 5}
	fmt.Println(s1)
	ModValue(s1)
	fmt.Println(s1)
}
