package middleware

import (
	"go-todo-api/domain"
	"go-todo-api/internal/config"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewAuth(userUseCase *usecase.UserUsecase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			userUseCase.Log.Warn("Failed to get user by token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Token is missing"})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			userUseCase.Log.Warn("Invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Invalid token format"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		userUseCase.Log.Debugf("Authorization: %s", token)

		jwt, _ := config.InitJwtService()
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			userUseCase.Log.Warnf("Failed to validate user token: %+v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
			return
		}

		getUserIdRequest := &domain.GetUserId{
			ID: claims.UserID,
		}

		user, err := userUseCase.GetUserID(ctx, getUserIdRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "User not found"})
			return
		}

		if user.Token == "" || user.Token != token {
			userUseCase.Log.Warnf("Failed to validate user token: %+v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Token invalid or missing"})
			return
		}

		ctx.Set("auth", user)

		ctx.Next()
	}
}

func GetUser(ctx *gin.Context) *entity.User {
	if auth, exists := ctx.Get("auth"); exists {
		return auth.(*entity.User)
	}
	return nil
}
