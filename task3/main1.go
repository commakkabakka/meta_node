package main

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

/*
	题目1：基本CRUD操作
		假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
	要求 ：
		编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
		编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
		编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
		编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
*/

type Student struct {
	gorm.Model
	Name  string
	Age   int
	Grade string
}

func Test31() {
	db, err := ConnectMysql()
	if err != nil {
		fmt.Println(err)
		return
	}

	db.AutoMigrate(&Student{})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"
	var student1 Student = Student{Name: "张三", Age: 20, Grade: "三年级"}
	err = gorm.G[Student](db).Create(ctx, &student1)
	if err != nil {
		fmt.Println(err)
	}

	// 查询 students 表中所有年龄大于 18 岁的学生信息。
	students, err := gorm.G[Student](db).Where("Age > ?", 18).Find(ctx)
	if err != nil {
		return
	}
	fmt.Println(students)

	// 将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	row, err := gorm.G[Student](db).Where("Name = ?", "张三").Update(ctx, "Grade", "四年级")
	if err != nil {
		return
	}
	fmt.Println(row)

	// 删除 students 表中年龄小于 15 岁的学生记录。
	row, err = gorm.G[Student](db).Where("Age < ?", 15).Delete(ctx)
	if err != nil {
		return
	}
	fmt.Println(row)
}
