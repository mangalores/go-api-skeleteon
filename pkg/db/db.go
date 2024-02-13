package db

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(c Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Berlin",
		c.Host,
		c.User,
		c.Password,
		c.DatabaseName,
		c.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.Default.LogMode(parseLogLevel(c.Logging))})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func parseLogLevel(logLevel string) logger.LogLevel {
	switch strings.ToLower(logLevel) {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "warning":
		return logger.Warn
	case "info":
		return logger.Info
	case "debug":
		return logger.Info
	case "trace":
		return logger.Info
	default:
		return logger.Silent
	}
}
