package main

import "fmt"

/*
	题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
*/

func PrintOdd(odd []int, ch chan int) {
	for _, value := range odd {
		if value%2 == 1 {
			fmt.Println("PrintOdd() : %d ", value)
		}
	}
	ch <- 1
}

func PrintEven(even []int, ch chan int) {
	for _, value := range even {
		if value%2 == 0 {
			fmt.Println("PrintEven() : %d ", value)
		}
	}
	ch <- 1
}

func Test23() {
	var ch chan int = make(chan int)

	var num []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	go PrintOdd(num, ch)
	go PrintEven(num, ch)

	// wait all goroutines
	for i := 0; i < 2; i++ {
		<-ch
	}
}
