package main

import (
	"fmt"
	"math"
)

/*
	题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
			在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
*/

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (a Rectangle) Area() {
	fmt.Println("Rectangle Area is ", a.Width*a.Height)
}

func (a Rectangle) Perimeter() {
	fmt.Println("Rectangle Perimeter is ", 2*a.Width+2*a.Height)
}

type Circle struct {
	R float64
}

func (a Circle) Area() {
	fmt.Println("Circle Area is ", math.Pi*a.R*a.R)
}

func (a Circle) Perimeter() {
	fmt.Println("Circle Perimeter is ", 2*math.Pi*a.R)
}

func Test25() {
	var circle Shape = Circle{2}
	var rectangle Shape = Rectangle{5, 3}

	circle.Area()
	circle.Perimeter()

	rectangle.Area()
	rectangle.Perimeter()
}
