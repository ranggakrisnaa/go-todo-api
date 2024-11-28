package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewGormConfig(cfg *GormConfig) *GormConfig {
	return &GormConfig{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASS"),
		DBName:   os.Getenv("DATABASE_NAME"),
	}
}

func (config *GormConfig) NewGormConnection() *gorm.DB {
	var dsn string
	if config.Password == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.DBName)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.Password, config.DBName)
	}
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from gorm.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return dbConn
}

func (config *GormConfig) RunMigrateDB() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	m, err := migrate.New(
		"file://db/migrations",
		dsn)
	if err != nil {
		log.Fatalf("migration initialization error: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Println("no new migrations to apply")
	}
}
