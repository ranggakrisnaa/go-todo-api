package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

func InitPostgreDB() (*gorm.DB, *GormConfig, error) {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	if host == "" || port == "" || user == "" || dbName == "" {
		return nil, nil, fmt.Errorf("database configuration is missing")
	}

	dbConfig := NewGormConfig(&GormConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	})
	return dbConfig.NewGormConnection(), dbConfig, nil
}

func InitJwtService() (*JwtConfig, error) {
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	jwtExpStr := os.Getenv("JWT_EXPIRATION_TIME")

	if jwtKey == "" || jwtExpStr == "" {
		return nil, fmt.Errorf("JWT configuration is incomplete: JWT_SECRET_KEY or JWT_EXPIRATION_TIME is missing")
	}

	jwtExp, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		return nil, err
	}

	jwtService, err := NewJwtConfig(&JwtConfig{
		JwtKey: jwtKey,
		JwtExp: jwtExp,
	})
	if err != nil {
		return nil, err
	}

	return jwtService, nil
}

func InitRedis() (*redis.Pool, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	maxActive, err := strconv.Atoi(os.Getenv("REDIS_MAX_ACTIVE"))
	if err != nil {
		return nil, err
	}
	maxIdle, err := strconv.Atoi(os.Getenv("REDIS_MAX_IDLE"))
	if err != nil {
		return nil, err
	}
	idleTimeout, err := strconv.Atoi(os.Getenv("REDIS_IDLE_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	if redisHost == "" || redisPort == "" {
		return nil, fmt.Errorf("missing configuration for redis")
	}

	return NewRedisPool(&RedisConfig{
		Host:        redisHost,
		Port:        redisPort,
		MaxActive:   maxActive,
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
	}), nil
}

func InitMailer() (*MailerConfig, error) {
	smtpHost := os.Getenv("CONFIG_SMTP_HOST")
	smtpPortStr := os.Getenv("CONFIG_SMTP_PORT")
	senderName := os.Getenv("CONFIG_SMTP_SENDER")
	smtpAuthEmail := os.Getenv("CONFIG_AUTH_EMAIL")
	smtpAuthPassword := os.Getenv("CONFIG_AUTH_PASSWORD")
	if smtpHost == "" || smtpPortStr == "" || smtpAuthEmail == "" || smtpAuthPassword == "" || senderName == "" {
		return nil, fmt.Errorf("SMTP configuration is missing")
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return nil, err
	}

	return NewMailerConfig(&MailerConfig{
		SmtpHost:         smtpHost,
		SmtpPort:         smtpPort,
		SenderMailName:   senderName,
		SmtpAuthEmail:    smtpAuthEmail,
		SmtpAuthPassword: smtpAuthPassword,
	}), nil
}
