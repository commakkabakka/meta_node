package main

import (
	"blog/middleware"
	"blog/mysql"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	err := mysql.InitGlobalDB()
	if err != nil {
		log.Fatalln(err)
	}
	// 创建数据库表
	//err = mysql.GetDBManager().CreateTable()
	//if err != nil {
	//	log.Fatalln(err)
	//}

	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	router.POST("/test", Test)

	router.POST("/register", mysql.Register)
	router.POST("/login", mysql.Login)

	blog := router.Group("/blog")
	blog.Use(middleware.JWTAuth())
	blog.POST("/create", mysql.BlogCreate)
	blog.POST("/list/:userid", mysql.BlogList)
	blog.POST("/detail/:blogid", mysql.BlogDetail)
	blog.POST("/update/:blogid", mysql.BlogUpdate)
	blog.POST("/delete/:blogid", mysql.BlogDelete)

	comment := router.Group("/comment")
	comment.Use(middleware.JWTAuth())
	comment.POST("/create/:blogid", mysql.CommentCreate)
	comment.POST("/list/:blogid", mysql.CommentList)

	router.Run() // 默认监听 0.0.0.0:8080
}

func Test(c *gin.Context) {
	panic("Something went wrong.")
}
