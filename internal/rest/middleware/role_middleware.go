package middleware

import (
	"go-todo-api/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequiredRole struct {
	User  string
	Admin string
}

func NewRequiredRole() *RequiredRole {
	return &RequiredRole{
		User:  "user",
		Admin: "admin",
	}
}

func (requiredRole *RequiredRole) RoleCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("auth")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized"})
			c.Abort()
			return
		}

		authUser, ok := user.(*entity.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized"})
			c.Abort()
			return
		}

		if authUser.Role != requiredRole.Admin {
			c.JSON(http.StatusForbidden, gin.H{"errors": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
