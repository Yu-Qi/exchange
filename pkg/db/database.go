package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

var database *gorm.DB
var initdatabaseOnce sync.Once

// Get database
func Get() *gorm.DB {
	initdatabaseOnce.Do(initialize)
	return database
}

// GetWith get database with context
func GetWith(ctx context.Context) *gorm.DB {
	initdatabaseOnce.Do(initialize)
	return database.WithContext(ctx)
}

// Init database
func Init() {
	initdatabaseOnce.Do(initialize)
}

func initialize() {
	var err error
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	options := os.Getenv("MYSQL_OPTIONS")
	databaseName := os.Getenv("MYSQL_DATABASE")
	slowThreshold := time.Second * 5
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", username, password, host, port, databaseName, options)

	database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             slowThreshold,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		// TODO: how to format error message for gcp error reporting
		log.Fatal(err)
	}
	if err = database.Use(tracing.NewPlugin()); err != nil {
		log.Fatal(err)
	}

}
