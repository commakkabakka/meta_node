package main

import (
	"fmt"

	"gorm.io/gorm"
)

/*
	题目2：实现类型安全映射
		假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
	要求 ：
		定义一个 Book 结构体，包含与 books 表对应的字段。
		编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/

type Book struct {
	gorm.Model
	Title  string
	Author string
	Price  float64
}

type BookInfo struct {
	Title  string  `db:"title";type"string"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func Test34() {
	db, err := ConnectMysql()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db)

	sqlx, err := ConnectMysql_sqlx()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sqlx)

	// 创建表
	db.AutoMigrate(&Book{})

	//var books []Book = []Book{
	//	{
	//		Title:  "AAA",
	//		Author: "ZS",
	//		Price:  52,
	//	},
	//	{
	//		Title:  "BBB",
	//		Author: "ZS",
	//		Price:  65,
	//	},
	//	{
	//		Title:  "CCC",
	//		Author: "ZS",
	//		Price:  37,
	//	},
	//}
	//db.Create(&books)

	var books_info []BookInfo
	err = sqlx.Select(&books_info, "SELECT title, author, price FROM books WHERE price > ?", 50)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(books_info)
}
