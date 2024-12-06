package main

import (
	"go-todo-api/internal"
	"go-todo-api/internal/config"
	"go-todo-api/internal/rest/middleware"
	"go-todo-api/internal/util"
	"go-todo-api/internal/workers"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
	_, db, _ := config.InitPostgreDB()
	db.RunMigrateDB()
}

func main() {
	db, _, err := config.InitPostgreDB()
	if err != nil {
		logrus.Fatalf("Failed to initialize mailer config: %v", err)
	}
	jwtService, err := config.InitJwtService()
	if err != nil {
		logrus.Fatalf("Failed to initialize mailer config: %v", err)
	}
	mailerConfig, err := config.InitMailer()
	if err != nil {
		logrus.Fatalf("Failed to initialize mailer config: %v", err)
	}
	redisPool, err := config.InitRedis()
	if err != nil {
		logrus.Fatalf("Failed to initialize redis pool config: %v", err)
	}

	workerPool := work.NewWorkerPool(workers.MailWorker{}, 10, "todo_queue", redisPool)
	mailWorker := workers.NewMailWorker(config.NewLogger(), mailerConfig)
	workerPool.Job("send_email", mailWorker.SendEmail)
	enqueuer := work.NewEnqueuer("todo_queue", redisPool)

	workerPool.Start()
	defer workerPool.Stop()

	r := gin.New()

	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		logrus.Fatalf("failed to parse timeout, using default timeout: %v", err)
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	r.Use(middleware.SetRequestContextWithTimeout(timeoutContext))
	r.Use(gin.LoggerWithFormatter(util.CustomLogFormatter))
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"success":    true,
			"statusCode": http.StatusOK,
			"message":    "Success get ping from server",
		})
	})

	internal.Bootstrap(&internal.BootstrapConfig{
		DB:         db,
		Log:        config.NewLogger(),
		Route:      r,
		JwtService: jwtService,
		Enqueurer:  enqueuer,
	})

	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}

	r.Run(address)
}
