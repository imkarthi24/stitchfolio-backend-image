package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConnectionParams struct {
	Host     string
	Port     int
	Username string
	DBName   string
	Password string
	SSLMode  string
	Schema   string
}

func ProvideDatabase(connectionParams DatabaseConnectionParams) (*gorm.DB, error) {

	host := connectionParams.Host
	port := connectionParams.Port
	userName := connectionParams.Username
	dbname := connectionParams.DBName
	password := connectionParams.Password
	schema := connectionParams.Schema
	sslMode := connectionParams.SSLMode

	args := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s search_path=%s", host, port, userName, dbname, password, sslMode, schema)

	connection, err := gorm.Open(postgres.Open(args), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Nanosecond, // Slow SQL threshold
				LogLevel:                  logger.Info,     // Log level (Silent, Error, Warn, Info)
				IgnoreRecordNotFoundError: false,           // Don't ignore ErrRecordNotFound error
				Colorful:                  true,            // Enable color
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

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error pinging Database: %v", err)
	}

	fmt.Println("Connected to database")
	return connection, nil
}
