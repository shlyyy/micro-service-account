package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwtutil "github.com/shlyyy/micro-service-account/pkg/jwt"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		tokenStr := c.GetHeader("Token")
		if tokenStr == "" {
			c.JSON(401, gin.H{"error": "缺少Token header"})
			c.Abort()
			return
		}

		claims, err := jwtutil.ParseToken(tokenStr)
		if err != nil {
			if err.Error() == jwtutil.ErrTokenExpired.Error() {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": jwtutil.ErrTokenExpired,
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "认证失败，需要登录",
			})
			c.Abort()
			return
		}

		// 将账户信息存到上下文，方便 handler 使用
		c.Set("claims", claims)
		c.Next()
	}
}
