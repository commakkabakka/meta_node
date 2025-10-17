package main

import "fmt"

/*
	题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。
			为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
*/

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

func (e Employee) PrintInfo() {
	fmt.Printf("%s is %d years old, his id is %d 。\n", e.Name, e.Age, e.EmployeeID)
}

func Test26() {
	var emp Employee = Employee{Person: Person{Name: "Hello", Age: 17}, EmployeeID: 123}
	emp.PrintInfo()
}
