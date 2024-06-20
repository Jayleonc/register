package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"time"
)

// RecoveryWithLogger 是一个自定义的恢复中间件，用于拦截 panic 并记录错误信息
func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 打印错误信息和堆栈
				fmt.Printf("[Recovery] %s panic recovered:\n%s\n%s\n",
					timeFormat(time.Now()), r, debug.Stack())

				// 可以在这里打日志

				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "服务器内部错误，请稍后再试",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// 时间格式化函数
func timeFormat(t time.Time) string {
	return t.Format("2006/01/02 - 15:04:05")
}
