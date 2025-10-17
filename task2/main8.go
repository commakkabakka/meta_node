package main

import (
	"fmt"
	"sync"
)

/*
	题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
*/

func Task28() {
	var wg sync.WaitGroup // 用于等待 goroutine 结束
	wg.Add(2)

	var ch = make(chan int, 100)

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
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
