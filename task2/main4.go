package main

import (
	"fmt"
	"time"
)

/*
	题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
*/

type TaskFunc func(task Task)

type Task struct {
	TaskFunc TaskFunc
	Args     []interface{}
	Ch       chan []interface{}
}

func RunTasks(tasks []Task) {
	for index := range tasks {
		go tasks[index].TaskFunc(tasks[index])
	}
}

func Test24() {
	tasks := []Task{
		{
			TaskFunc: func(task Task) {
				start := time.Now()
				defer func() {
					fmt.Println("任务执行时间 ：", time.Since(start))
				}()
				a := task.Args[0].(int)
				b := task.Args[1].(int)
				time.Sleep(2 * time.Second)
				task.Ch <- []interface{}{a + b} // return one value
			},
			Args: []interface{}{3, 7},
			Ch:   make(chan []interface{}),
		},
		{
			TaskFunc: func(task Task) {
				start := time.Now()
				defer func() {
					fmt.Println("任务执行时间 ：", time.Since(start))
				}()
				name := task.Args[0].(string)
				time.Sleep(4 * time.Second)
				task.Ch <- []interface{}{fmt.Sprintf("Hello, %s.", name), 12345} // return two value
			},
			Args: []interface{}{"World"},
			Ch:   make(chan []interface{}),
		},
		{
			TaskFunc: func(task Task) {
				start := time.Now()
				defer func() {
					fmt.Println("任务执行时间 ：", time.Since(start))
				}()
				time.Sleep(6 * time.Second)
				task.Ch <- []interface{}{"ok", 100, true} // 返回三个值
			},
			Args: nil,
			Ch:   make(chan []interface{}),
		},
	}

	RunTasks(tasks)

	for i := 0; i < len(tasks); i++ {
		select {
		case rets := <-tasks[0].Ch:
			{
				fmt.Println("任务1执行完成。结果 ： ", rets)
			}
		case rets := <-tasks[1].Ch:
			{
				fmt.Println("任务2执行完成。结果 ： ", rets)
			}
		case rets := <-tasks[2].Ch:
			{
				fmt.Println("任务3执行完成。结果 ： ", rets)
			}
		}
	}

}
