package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/inflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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

func init() {
	// Clear default irregular rules to prevent unwanted pluralization
	inflection.SetIrregular(inflection.IrregularSlice{})
}

// ProvideDatabase initializes and returns a gorm DB connection
func ProvideDatabase(connectionParams DatabaseConnectionParams) (*gorm.DB, error) {

	host := connectionParams.Host
	port := connectionParams.Port
	userName := connectionParams.Username
	dbname := connectionParams.DBName
	password := connectionParams.Password
	dbSchema := connectionParams.Schema
	sslMode := connectionParams.SSLMode

	args := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s search_path=%s", host, port, userName, dbname, password, sslMode, dbSchema)

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
		NamingStrategy: CustomStrategy{
			NamingStrategy: schema.NamingStrategy{
				NoLowerCase:   false, //  keep columns snake_case and for table we resolve using custom strategy
				SingularTable: false,
			},
			Schema: dbSchema,
		},
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

type CustomStrategy struct {
	schema.NamingStrategy
	Schema string
}

func (ns CustomStrategy) TableName(str string) string {
	// original = struct name (e.g. UserChannelDetail)
	table := str

	// pluralize if gorm needs it (you can disable if you want single)
	if !ns.SingularTable {
		table = inflection.Plural(str)
	}

	// attach schema prefix
	if ns.Schema != "" {
		return fmt.Sprintf(`%s.%s`, ns.Schema, table)
	}

	return table
}
