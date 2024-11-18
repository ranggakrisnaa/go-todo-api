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
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	defaultTimeout = 10
	defaultAddress = ":8080"
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
	logger := config.NewLogger()

	r := gin.New()

	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	r.Use(middleware.SetRequestContextWithTimeout(timeoutContext))
	r.Use(gin.LoggerWithFormatter(util.CustomLogFormatter))
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.GET("/ping", func(c *gin.Context) {
		util.SendSuccess(c, http.StatusOK, "Success ping the server", nil)
	})

	userRepo := postgresql.NewRepository(dbConn)
	userService := usecase.NewUserUseCase(userRepo, dbConn, logger)
	rest.NewUserHandler(r, userService, logger)

	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}

	r.Run(address)
}
