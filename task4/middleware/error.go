package middleware

import (
	"blog/error"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 捕获 panic
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered: %v, path: %s", r, c.FullPath())
				c.JSON(
					http.StatusInternalServerError,
					gin.H{
						"code":    http.StatusInternalServerError,
						"message": "internal server error",
					})
				c.Abort()
			}
		}()

		c.Next()

		// 统一处理业务错误
		if len(c.Errors) > 0 {
			// 记录日志
			for _, e := range c.Errors {
				log.Printf("handler error: %v", e.Err)
			}

			// 统一返回错误信息。
			firstErr := c.Errors[0].Err
			appErr, ok := firstErr.(*error.AppError)
			if ok {
				c.JSON(
					appErr.Code,
					gin.H{
						"code":  appErr.Code,
						"error": appErr.Message,
					})
			} else {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{
						"code":  http.StatusInternalServerError,
						"error": "Internal server error.",
					})
			}
		}

	}
}
