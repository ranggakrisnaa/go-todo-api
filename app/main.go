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
	"github.com/gomodule/redigo/redis"
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
	db := config.InitPostgreDB()
	logger := config.InitLogger()
	jwtService := config.InitJwtService(logger)

	mailerConfig, err := config.NewMailerConfig()
	if err != nil {
		log.Fatalf("Failed to initialize mailer config: %v", err)
	}
	redisPool := &redis.Pool{
		MaxActive:   50,
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	workerPool := work.NewWorkerPool(workers.MailWorker{}, 10, "todo_queue", redisPool)

	mailWorker := workers.NewMailWorker(logger, mailerConfig)
	workerPool.Job("send_email", mailWorker.SendEmail)
	enqueuer := work.NewEnqueuer("todo_queue", redisPool)
	workerPool.Start()
	defer workerPool.Stop()

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
		Log:        logger,
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
