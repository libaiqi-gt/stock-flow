package middleware

import (
	"stock-flow/internal/pkg/response"
	"stock-flow/internal/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(c, response.CodeUnauthorized, "请求未携带Token，无权访问")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, response.CodeUnauthorized, "Token格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			response.Error(c, response.CodeUnauthorized, "Token无效或已过期")
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Error(c, response.CodeUnauthorized, "未获取到角色信息")
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		for _, role := range roles {
			if role == roleStr || roleStr == "Admin" { // Admin has all permissions
				c.Next()
				return
			}
		}

		response.Error(c, response.CodeForbidden, "权限不足")
		c.Abort()
	}
}
