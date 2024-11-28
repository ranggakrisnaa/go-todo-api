package internal

import (
	"go-todo-api/internal/config"
	"go-todo-api/internal/repository/postgresql"
	"go-todo-api/internal/rest"
	"go-todo-api/internal/rest/middleware"
	"go-todo-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB         *gorm.DB
	Route      *gin.Engine
	Log        *logrus.Logger
	JwtService *config.JwtConfig
}

func Bootstrap(config *BootstrapConfig) {
	userRepo := postgresql.NewUserRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, config.DB, config.Log, config.JwtService)
	authMiddleware := middleware.NewAuth(userUsecase)
	rest.NewUserHandler(config.Route, userUsecase, config.Log, authMiddleware)

	todoRepo := postgresql.NewTodoRepository(config.DB)
	todoUsecase := usecase.NewTodoUseCase(todoRepo, config.DB, config.Log, config.JwtService)
	rest.NewTodoHandler(config.Route, todoUsecase, config.Log)

}
