package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

/*
	题目1：使用SQL扩展库进行查询
		假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	要求 ：
		编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
		编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/

// Employee 结构体映射 employees 表
type Employee struct {
	gorm.Model
	Name       string
	Department string
	Salary     float64
}

func Test33() {
	db, err := ConnectMysql()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlx, err := ConnectMysql_sqlx()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建表
	db.AutoMigrate(&Employee{})

	//var employees []Employee = []Employee{
	//	{
	//		Name:       "AAA",
	//		Department: "市场部",
	//		Salary:     8000,
	//	},
	//	{
	//		Name:       "BBB",
	//		Department: "技术部",
	//		Salary:     9000,
	//	},
	//	{
	//		Name:       "CCC",
	//		Department: "财务部",
	//		Salary:     10000,
	//	},
	//}
	//db.Create(&employees)

	// 注意 MySQL 大小写敏感
	var employees []Employee
	err = sqlx.Select(&employees, "SELECT id, name, department, salary FROM employees WHERE DEPARTMENT = ?", "技术部")
	if err != nil {
		log.Fatalln("查询失败1:", err)
	}
	fmt.Println(employees)

	var employee Employee
	err = sqlx.Get(&employee, "SELECT id, name, department, salary FROM employees WHERE DEPARTMENT = ? ORDER BY  salary DESC LIMIT 1", "技术部")
	if err != nil {
		log.Fatalln("查询失败2:", err)
	}
	fmt.Println(employee)
}
