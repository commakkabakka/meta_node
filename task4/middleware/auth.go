package middleware

import (
	"blog/mysql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}
		// 去掉 Bearer 前缀
		parts := strings.Fields(tokenString)
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token : " + tokenString})
			c.Abort()
			return
		}
		tokenString = parts[1]

		// 解析 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 强制校验签名方法为 HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte("my_secret_key"), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token : " + err.Error()})
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token : valid."})
			c.Abort()
			return
		}

		// 转换为 MapClaims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token : format err."})
			c.Abort()
			return
		}

		// 检查过期时间 - 当 JSON 被解码到 interface{} 时，Go 默认会把所有数字 解码为 float64。
		exp, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"ok": "Invalid token : type err."})
			c.Abort()
			return
		}

		if time.Now().Unix() > (int64)(exp) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token : expired."})
			c.Abort()
			return
		} else {
			userid := (uint)(claims["userid"].(float64))
			username := claims["username"].(string)

			// 找不到这个用户
			var user mysql.User
			err := mysql.Ins.DB.Where("id = ?", userid).First(&user)
			if err.Error != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error})
				c.Abort()
				return
			}

			c.Set("BLOG_userid", userid)
			c.Set("BLOG_username", username)

			log.Println("---------------------- user_id = ?", userid)

			//c.JSON(http.StatusOK, gin.H{"ok": "验证通过 ：" + username})
			c.Next()
		}
	}
}
