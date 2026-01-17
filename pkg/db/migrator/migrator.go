package migrator

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"
)

type migrator struct {
	dB *gorm.DB
}

type Migrate interface {
	Migrate(entities []interface{}, checkErr func(err error))
	GenerateAlterMigration(entities []interface{}, migrationName string) error
}

func NewMigrator(db *gorm.DB) Migrate {
	return &migrator{
		dB: db,
	}
}

// Migrate generates an initial migration file based on the current state of the models
func (m *migrator) Migrate(entities []interface{}, checkErr func(err error)) {

	fmt.Println("Generating SQL migration...")

	sqlStatements, err := generateSQL(m.dB, entities)
	if err != nil {
		checkErr(err)
		return
	}

	// Create migrations directory if it doesn't exist
	os.MkdirAll("migrations", 0755)

	// Write to file
	file, err := os.Create("migrations/001_initial.sql")
	if err != nil {
		fmt.Printf("❌ Error creating file: %v\n", err)
		checkErr(err)
		return
	}
	defer file.Close()

	// Write header
	file.WriteString("-- Migration: Initial Schema\n")
	file.WriteString("-- Generated automatically from GORM models\n\n")

	// Separate table creation from constraints
	file.WriteString("-- ====================================\n")
	file.WriteString("-- CREATE TABLES\n")
	file.WriteString("-- ====================================\n\n")

	createTableCount := 0
	indexCount := 0
	fkCount := 0

	for _, sql := range sqlStatements {
		if strings.HasPrefix(sql, "CREATE TABLE") {
			file.WriteString(sql + ";\n\n")
			createTableCount++
		} else if strings.HasPrefix(sql, "CREATE") {
			// This is an index
			if indexCount == 0 {
				file.WriteString("\n-- ====================================\n")
				file.WriteString("-- CREATE INDEXES\n")
				file.WriteString("-- ====================================\n\n")
			}
			file.WriteString(sql + ";\n\n")
			indexCount++
		} else if strings.HasPrefix(sql, "ALTER TABLE") {
			// This is a foreign key
			if fkCount == 0 {
				file.WriteString("\n-- ====================================\n")
				file.WriteString("-- ADD FOREIGN KEY CONSTRAINTS\n")
				file.WriteString("-- ====================================\n\n")
			}
			file.WriteString(sql + ";\n\n")
			fkCount++
		}
	}

	file.WriteString("\n-- Migration completed successfully\n")
	file.WriteString(fmt.Sprintf("-- Tables created: %d\n", createTableCount))
	file.WriteString(fmt.Sprintf("-- Indexes created: %d\n", indexCount))
	file.WriteString(fmt.Sprintf("-- Foreign keys added: %d\n", fkCount))

	fmt.Printf("✓ Migration file created: migrations/001_initial.sql\n")
	fmt.Printf("  - Tables: %d\n", createTableCount)
	fmt.Printf("  - Indexes: %d\n", indexCount)
	fmt.Printf("  - Foreign Keys: %d\n", fkCount)
}

// GenerateAlterMigration generates an alter migration file based on the current state of the models
func (m *migrator) GenerateAlterMigration(entities []interface{}, migrationName string) error {
	return m.generateMigration(entities, migrationName)
}
