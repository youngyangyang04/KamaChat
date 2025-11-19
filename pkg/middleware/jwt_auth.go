package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwtutil "kama_chat_server/pkg/util/jwt"
	"kama_chat_server/pkg/zlog"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请求头中缺少 Authorization",
			})
			c.Abort()
			return
		}

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

		tokenString := parts[1]

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
