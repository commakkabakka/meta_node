package main

import (
	"fmt"

	"gorm.io/gorm"
)

/*
	题目1：模型定义
		假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
	要求 ：
		使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
		编写Go代码，使用Gorm创建这些模型对应的数据库表。
*/

/*
	题目2：关联查询
		基于上述博客系统的模型定义。
	要求 ：
		编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
		编写Go代码，使用Gorm查询评论数量最多的文章信息。
*/

/*
	题目3：钩子函数
		继续使用博客系统的模型。
	要求 ：
		为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
		为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/

/*
	注意点：
		1. Create supports nested associations. Hooks will be triggered only for the first-level associations, nested deeper will not trigger hooks.
		2. When bulk deleting, hooks are called only once, and the model parameter passed into the hook will not have field values populated.
*/

type User struct {
	gorm.Model
	Name      string
	Posts     []Post
	PostCount uint
}

type Post struct {
	gorm.Model
	Title        string
	Content      string
	UserID       uint
	Comments     []Comment
	CommentState string
}

func (post *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).Where("id = ?", post.UserID).Update("post_count", gorm.Expr("post_count+1")).Error
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

func (comment *Comment) AfterCreate(tx *gorm.DB) error {
	var count int64
	// 查询该文章的剩余评论数
	err := tx.Model(&Comment{}).Where("post_id = ?", comment.PostID).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 1 {
		return tx.Model(&Post{}).Where("id = ?", comment.PostID).Update("comment_state", "有评论").Error
	}

	return nil
}

func (comment *Comment) AfterDelete(tx *gorm.DB) error {
	var count int64
	err := tx.Model(&Comment{}).Where("post_id = ?", comment.PostID).Count(&count).Error
	if err != nil {
		return err
	}

	fmt.Println(count)

	if count == 0 {
		return tx.Model(&Post{}).Where("id = ?", comment.PostID).Update("comment_state", "无评论").Error
	}

	return nil
}

type Result struct {
	PostID  uint
	Title   string
	Content string
	Count   int64
}

func Test35() {
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

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Comment{})

}

func Test36() {
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

	var users []User = []User{
		{
			Name: "AAA",
			Posts: []Post{
				{
					Title:   "A-1",
					Content: "你哈山豆根立卡就是都发了噶事g.111111",
					Comments: []Comment{
						{
							Content: "可以",
						},
						{
							Content: "不错",
						},
						{
							Content: "非常棒",
						},
					},
					CommentState: "-",
				},
				{
					Title:        "A-2",
					Content:      "你哈山豆根立卡就是都发了噶事g.222222",
					Comments:     []Comment{},
					CommentState: "-",
				},
				{
					Title:   "A-3",
					Content: "你哈山豆根立卡就是都发了噶事g.3333333",
					Comments: []Comment{
						{
							Content: "可以333",
						},
						{
							Content: "不错333",
						},
						{
							Content: "非常棒2333",
						},
					},
					CommentState: "-",
				},
			},
			PostCount: 0,
		},
	}

	db.Create(&users)

	// 问题一 ：编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	var user User
	db.Preload("Posts.Comments").First(&user, "name = ?", "AAA")
	fmt.Println(user)

	// 问题二 ：编写Go代码，使用Gorm查询评论数量最多的文章信息。
	var ret Result
	err = db.Model(&Post{}).
		Select("posts.id as post_id, posts.title, posts.content, count(comments.id) as count").
		Joins("join comments on posts.id = comments.post_id").
		Group("posts.id").
		Order("count DESC").
		Limit(1).
		Scan(&ret).Error
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ret)
}

func Test37() {
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

	// 测试：为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	//var post Post = Post{
	//	Title:   "12345",
	//	Content: "来看哈森岛帆高哈市dg",
	//	UserID:  1,
	//}
	//db.Create(&post)

	// 测试：为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
	//var comment Comment = Comment{
	//	Content: "不错",
	//	PostID:  2,
	//}
	//db.Create(&comment)

	// 为了使钩子函数能够获取 PostID ，必须在这里赋值。
	var comment Comment
	comment.ID = 8
	comment.PostID = 2
	db.Debug().Delete(&comment)

	// 推荐方式 ：先查询再删除。
	//var comment Comment
	//if err := db.First(&comment, 8).Error; err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//db.Delete(&comment) // 单条删除，钩子中 comment.PostID 有值
}
