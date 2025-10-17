package mysql

import (
	"blog/error"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

func CommentCreate(c *gin.Context) {
	blogID := c.Param("blogid")

	var post Post
	err := Ins.DB.Debug().Where("id = ?", blogID).First(&post).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	var comment Comment
	err = c.ShouldBindJSON(&comment)
	if err != nil {
		_ = c.Error(error.BadRequest)
		c.Abort()
		return
	}

	userid := c.GetUint("BLOG_userid")

	comment.PostID = post.ID
	comment.UserID = userid

	err = Ins.DB.Debug().Create(&comment).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Create comment success.",
	})
}

func CommentList(c *gin.Context) {
	blogID := c.Param("blogid")

	var comments []Comment
	err := Ins.DB.Debug().Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).Preload("Post", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, title, content")
	}).Where("post_id = ?", blogID).Find(&comments).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"message":  "List comments success.",
		"comments": comments,
	})
}
