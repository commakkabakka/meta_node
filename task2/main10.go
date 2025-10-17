package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
	题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
*/

func Task210() {
	var wg sync.WaitGroup
	var num atomic.Int32 = atomic.Int32{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				num.Add(1)
			}
		}()
	}

	wg.Wait()
	fmt.Println(num.Load())
}
