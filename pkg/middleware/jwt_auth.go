package middleware

import (
	"net/http"
	"strings"

	jwtutil "kama_chat_server/pkg/util/jwt"
	"kama_chat_server/pkg/zlog"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 优先从请求头获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// 解析 token（格式：Bearer <token>）
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Authorization 格式错误，应为 Bearer <token>",
				})
				c.Abort()
				return
			}
			tokenString = parts[1]
		} else {
			// 如果请求头中没有 token，尝试从 query 参数获取（用于 WebSocket）
			tokenString = c.Query("token")
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "缺少认证 token",
				})
				c.Abort()
				return
			}
		}

		// 验证 token
		claims, err := jwtutil.ParseToken(tokenString)
		if err != nil {
			zlog.Error("Token 验证失败: " + err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文，供后续处理使用
		c.Set("uuid", claims.UUID)
		c.Set("account", claims.Account)
		c.Set("nickname", claims.Nickname)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}
