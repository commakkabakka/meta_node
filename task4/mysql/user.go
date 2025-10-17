package mysql

import (
	"blog/error"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

// 用户注册
func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		_ = c.Error(error.BadRequest)
		c.Abort()
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		_ = c.Error(error.ErrInternalServer)
		c.Abort()
		return
	}
	user.Password = string(hashedPassword)

	if err := Ins.DB.Create(&user).Error; err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Register user success.",
	})
}

// 用户登陆
func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		_ = c.Error(error.BadRequest)
		c.Abort()
		return
	}

	var storedUser User
	if err := Ins.DB.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		_ = c.Error(error.ErrDatabase)
		c.Abort()
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		_ = c.Error(error.ErrUnauthorized)
		c.Abort()
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":   storedUser.ID,
		"username": storedUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		_ = c.Error(error.ErrUnauthorized)
		c.Abort()
		return
	}
	// 剩下的逻辑...

	// 方式一
	//c.SetCookie("access_token", tokenString, 3600*24, "/", "example.com", true, true)

	// 方式二: 使用 JSON 方便 POSTMAN 调试。
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Login success.",
		"token":   tokenString,
	})
}
