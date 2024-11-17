package main

import (
	"go-todo-api/internal/config"
	"go-todo-api/internal/repository/postgresql"
	"go-todo-api/internal/rest"
	"go-todo-api/internal/rest/middleware"
	"go-todo-api/internal/usecase"
	"go-todo-api/internal/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.RunMigrateDB()
}

func main() {
	configDB := config.NewPostgresConfig()
	dbConn := configDB.NewPostgresConnection()

	r := gin.New()

	r.Use(gin.LoggerWithFormatter(util.CustomLogFormatter))
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.GET("/ping", func(c *gin.Context) {
		util.SendSuccess(c, http.StatusOK, "Success ping the server", nil)
	})

	userRepo := postgresql.NewRepository(dbConn)
	userService := usecase.NewUserUseCase(userRepo)
	rest.NewUserHandler(r, userService)

	r.Run(":8080")
}
