package main

import (
	"fmt"
	"sync"
)

/*
	题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
*/

func Task29() {
	var wg sync.WaitGroup
	var mu sync.Mutex = sync.Mutex{}
	var num int = 0

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				func() {
					defer mu.Unlock()
					num++
				}()
			}
		}()
	}

	wg.Wait()
	fmt.Println(num)
}
