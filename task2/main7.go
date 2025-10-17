package main

import (
	"fmt"
	"sync"
)

/*
	题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
*/

func Task27() {
	var wg sync.WaitGroup // 用于等待 goroutine 结束
	wg.Add(2)

	var ch = make(chan int)

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		defer wg.Done()
		for {
			value, ok := <-ch
			if !ok {
				break
			}
			fmt.Println("Read : ", value)
		}
	}()

	wg.Wait()
}
