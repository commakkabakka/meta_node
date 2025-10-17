package mysql

import (
	"blog/error"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
}

func BlogCreate(c *gin.Context) {
	var post Post
	err := c.ShouldBindJSON(&post)
	if err != nil {
		_ = c.Error(error.BadRequest)
		c.Abort()
		return
	}

	userid := c.GetUint("BLOG_userid")
	post.UserID = userid
	err = Ins.DB.Debug().Create(&post).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Create post success.",
	})
}

func BlogList(c *gin.Context) {
	userID := c.Param("userid")

	var posts []Post
	err := Ins.DB.Debug().Where("user_id = ?", userID).Find(&posts).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "List posts success.",
		"posts":   posts,
	})
}

func BlogDetail(c *gin.Context) {
	blogID := c.Param("blogid")

	var post Post
	err := Ins.DB.Debug().Where("id = ?", blogID).First(&post).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "List post success.",
		"post":    post,
	})
}

func BlogUpdate(c *gin.Context) {
	blogID := c.Param("blogid")
	var postNew Post
	err := c.ShouldBindJSON(&postNew)
	if err != nil {
		_ = c.Error(error.BadRequest)
		c.Abort()
		return
	}
	userid := c.GetUint("BLOG_userid")

	var postOld Post
	err = Ins.DB.Debug().Where("id = ?", blogID).First(&postOld).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	if postOld.UserID != userid {
		_ = c.Error(error.ErrUnauthorized)
		c.Abort()
		return
	}
	postOld.Title = postNew.Title
	postOld.Content = postNew.Content
	err = Ins.DB.Save(&postOld).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Update post success.",
	})
}

func BlogDelete(c *gin.Context) {
	blogID := c.Param("blogid")
	userid := c.GetUint("BLOG_userid")

	var postOld Post
	err := Ins.DB.Debug().Where("id = ?", blogID).First(&postOld).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	if postOld.UserID != userid {
		_ = c.Error(error.ErrUnauthorized)
		c.Abort()
		return
	}
	log.Println("postOld.UserID != userid : ? , ?", postOld.UserID, userid)

	err = Ins.DB.Debug().Delete(&postOld).Error
	if err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Delete post success.",
	})
}
