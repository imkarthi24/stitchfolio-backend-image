package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/imkarthi24/sf-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ProvideDatabase(config config.DatabaseConfig) (*gorm.DB, error) {

	host := config.Host
	port := config.Port
	userName := config.Username
	dbname := config.DBName
	password := config.Password

	args := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s options='-c search_path=%s'", host, port, userName, dbname, password, "prefer", config.Schema)

	connection, err := gorm.Open(postgres.Open(args), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Millisecond, // Slow SQL threshold
				LogLevel:                  logger.Info,      // Log level (Silent, Error, Warn, Info)
				IgnoreRecordNotFoundError: false,            // Don't ignore ErrRecordNotFound error
				Colorful:                  true,             // Enable color
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("Error Connecting to Database : %v", err)
	}

	db, err := connection.DB()
	if err != nil {
		return nil, fmt.Errorf("Error Connecting to Database : %v", err)
	}

	if _ = connection.Exec(fmt.Sprintf("SET search_path TO %s", config.Schema)); err != nil {
		return nil, fmt.Errorf("failed to set search_path: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error pinging Database: %v", err)
	}

	fmt.Println("Connected to database")
	return connection, nil
}
