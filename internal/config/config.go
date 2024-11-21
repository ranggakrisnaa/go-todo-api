package config

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitPostgres() *gorm.DB {
	postgresqlConf := NewPostgresConfig()
	return postgresqlConf.NewPostgresConnection()
}

func InitLogger() *logrus.Logger {
	return NewLogger()
}

func InitJwtService(logger *logrus.Logger) *JwtConfig {
	jwtService, err := NewJwtConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize JwtConfig")
		return nil
	}
	return jwtService
}
