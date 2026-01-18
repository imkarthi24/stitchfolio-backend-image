package migrator

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type SchemaChange struct {
	Type      string // "ADD_COLUMN", "DROP_COLUMN", "MODIFY_COLUMN", "ADD_INDEX", "DROP_INDEX", "ADD_FK", "DROP_FK"
	TableName string
	SQL       string
}

// GenerateMigration compares current DB schema with models and generates ALTER scripts
func (sd *migrator) generateMigration(models []interface{}, migrationName string) error {
	var changes []SchemaChange

	for _, model := range models {
		stmt := &gorm.Statement{DB: sd.dB}
		if err := stmt.Parse(model); err != nil {
			return err
		}

		tableName := stmt.Schema.Table
		// Check if table exists
		if !sd.dB.Migrator().HasTable(model) {
			// Table doesn't exist - create it
			sql := buildCreateTableSQL(stmt.Schema)
			changes = append(changes, SchemaChange{
				Type:      "CREATE_TABLE",
				TableName: tableName,
				SQL:       sql,
			})

			// Add indexes
			for _, indexSQL := range buildIndexSQL(stmt.Schema) {
				changes = append(changes, SchemaChange{
					Type:      "ADD_INDEX",
					TableName: tableName,
					SQL:       indexSQL,
				})
			}

			// Add foreign keys
			for _, fkSQL := range buildForeignKeySQL(stmt.Schema) {
				changes = append(changes, SchemaChange{
					Type:      "ADD_FK",
					TableName: tableName,
					SQL:       fkSQL,
				})
			}
		} else {
			// Table exists - check for changes
			columnChanges := sd.detectColumnChanges(stmt.Schema)
			changes = append(changes, columnChanges...)

			indexChanges := sd.detectIndexChanges(stmt.Schema)
			changes = append(changes, indexChanges...)
		}
	}

	if len(changes) == 0 {
		fmt.Println("✓ No schema changes detected")
		return nil
	}

	// Generate migration file
	return sd.writeMigrationFile(changes, migrationName)
}

func (sd *migrator) detectColumnChanges(schema *schema.Schema) []SchemaChange {
	var changes []SchemaChange
	tableName := schema.Table

	// Get existing columns from database
	columnTypes, err := sd.dB.Migrator().ColumnTypes(tableName)
	if err != nil {
		return changes
	}

	existingColumns := make(map[string]bool)
	for _, col := range columnTypes {
		existingColumns[col.Name()] = true
	}

	// Check for new columns
	for _, field := range schema.Fields {
		if field.DBName == "" {
			continue
		}

		if !sd.dB.Migrator().HasColumn(tableName, field.DBName) {
			// Column doesn't exist - add it
			sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s",
				getTableNameWithQuotes(schema),
				buildColumnDefinition(field))
			changes = append(changes, SchemaChange{
				Type:      "ADD_COLUMN",
				TableName: tableName,
				SQL:       sql,
			})
		} else {
			// Column exists - check if it needs modification
			modifySQL := sd.detectColumnModification(tableName, field)
			if modifySQL != "" {
				changes = append(changes, SchemaChange{
					Type:      "MODIFY_COLUMN",
					TableName: tableName,
					SQL:       modifySQL,
				})
			}
		}
	}

	// Check for dropped columns
	modelColumns := make(map[string]bool)
	for _, field := range schema.Fields {
		if field.DBName != "" {
			modelColumns[field.DBName] = true
		}
	}

	for colName := range existingColumns {
		if !modelColumns[colName] {
			// Column exists in DB but not in model - should be dropped
			sql := fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s",
				tableName,
				colName)
			changes = append(changes, SchemaChange{
				Type:      "DROP_COLUMN",
				TableName: tableName,
				SQL:       sql,
			})
		}
	}

	return changes
}

func (sd *migrator) detectColumnModification(tableName string, field *schema.Field) string {
	// Get column type from database
	columnTypes, err := sd.dB.Migrator().ColumnTypes(tableName)
	if err != nil {
		return ""
	}

	var dbColumn *gorm.ColumnType
	for _, col := range columnTypes {
		if col.Name() == field.DBName {
			dbColumn = &col
			break
		}
	}

	if dbColumn == nil {
		return ""
	}

	var alterations []string

	// Check data type change
	expectedType := getPostgresType(field)
	dbType := (*dbColumn).DatabaseTypeName()

	// Normalize types for comparison
	if !sd.typesMatch(expectedType, dbType) {
		alterations = append(alterations,
			fmt.Sprintf("ALTER COLUMN %s TYPE %s", field.DBName, expectedType))
	}

	// Check nullable change
	nullable, ok := (*dbColumn).Nullable()
	if ok {
		if field.NotNull && nullable {
			alterations = append(alterations,
				fmt.Sprintf("ALTER COLUMN %s SET NOT NULL", field.DBName))
		} else if !field.NotNull && !nullable {
			alterations = append(alterations,
				fmt.Sprintf("ALTER COLUMN %s DROP NOT NULL", field.DBName))
		}
	}

	defaultDBValue, _ := (*dbColumn).DefaultValue()
	// Check default value change
	if field.DefaultValue != "" && defaultDBValue != field.DefaultValue {
		if field.DataType == "text" && !strings.HasPrefix(field.DefaultValue, "'") {
			//Add single quotes for text default values
			field.DefaultValue = fmt.Sprintf("'%s'", field.DefaultValue)
		}
		alterations = append(alterations,
			fmt.Sprintf("ALTER COLUMN %s SET DEFAULT %s", field.DBName, field.DefaultValue))
	}

	if len(alterations) == 0 {
		return ""
	}

	return fmt.Sprintf("ALTER TABLE %s %s", tableName, strings.Join(alterations, ", "))
}

func (sd *migrator) detectIndexChanges(schema *schema.Schema) []SchemaChange {
	var changes []SchemaChange
	tableName := getTableNameWithQuotes(schema)

	// Get existing indexes
	existingIndexes, err := sd.dB.Migrator().GetIndexes(tableName)
	if err != nil {
		return changes
	}

	existingIndexMap := make(map[string]bool)
	for _, idx := range existingIndexes {
		existingIndexMap[idx.Name()] = true
	}

	// Check for new indexes
	for _, index := range schema.ParseIndexes() {
		if !existingIndexMap[index.Name] {
			var fields []string
			for _, field := range index.Fields {
				fields = append(fields, field.DBName)
			}

			indexType := "INDEX"
			if index.Class == "UNIQUE" {
				indexType = "UNIQUE INDEX"
			}

			sql := fmt.Sprintf("CREATE %s IF NOT EXISTS %s ON %s (%s)",
				indexType,
				index.Name,
				tableName,
				strings.Join(fields, ", "))

			changes = append(changes, SchemaChange{
				Type:      "ADD_INDEX",
				TableName: tableName,
				SQL:       sql,
			})
		}
	}

	return changes
}

func (sd *migrator) writeMigrationFile(changes []SchemaChange, migrationName string) error {
	// Create migrations directory
	os.MkdirAll("migrations", 0755)

	// Generate timestamp

	filename := fmt.Sprintf("migrations/%s.sql", migrationName)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	file.WriteString(fmt.Sprintf("-- Migration: %s\n", migrationName))
	file.WriteString(fmt.Sprintf("-- Generated: %s\n\n", time.Now().Format(time.RFC3339)))

	// Group changes by type
	file.WriteString("-- ====================================\n")
	file.WriteString("-- UP Migration\n")
	file.WriteString("-- ====================================\n\n")

	// Create tables first
	for _, change := range changes {
		if change.Type == "CREATE_TABLE" {
			file.WriteString(fmt.Sprintf("-- Create table: %s\n", change.TableName))
			file.WriteString(change.SQL + ";\n\n")
		}
	}

	// Add columns
	for _, change := range changes {
		if change.Type == "ADD_COLUMN" {
			file.WriteString(fmt.Sprintf("-- Add column to %s\n", change.TableName))
			file.WriteString(change.SQL + ";\n\n")
		}
	}

	// Modify columns
	for _, change := range changes {
		if change.Type == "MODIFY_COLUMN" {
			file.WriteString(fmt.Sprintf("-- Modify column in %s\n", change.TableName))
			file.WriteString(change.SQL + ";\n\n")
		}
	}

	// Add indexes
	for _, change := range changes {
		if change.Type == "ADD_INDEX" {
			file.WriteString(fmt.Sprintf("-- Add index to %s\n", change.TableName))
			file.WriteString(change.SQL + ";\n\n")
		}
	}

	// Add foreign keys
	for _, change := range changes {
		if change.Type == "ADD_FK" {
			file.WriteString(fmt.Sprintf("-- Add foreign key to %s\n", change.TableName))
			file.WriteString(change.SQL + ";\n\n")
		}
	}

	// Drop columns (at the end)
	for _, change := range changes {
		if change.Type == "DROP_COLUMN" {
			file.WriteString(fmt.Sprintf("-- Drop column from %s\n", change.TableName))
			file.WriteString(change.SQL + ";\n\n")
		}
	}

	// Write DOWN migration
	file.WriteString("\n-- ====================================\n")
	file.WriteString("-- DOWN Migration (Rollback)\n")
	file.WriteString("-- ====================================\n\n")
	file.WriteString("-- TODO: Add rollback statements manually\n")

	fmt.Printf("✓ Migration file created: %s\n", filename)
	fmt.Printf("  Total changes: %d\n", len(changes))

	return nil
}

func (sd *migrator) typesMatch(expected, actual string) bool {
	// Normalize type names for comparison
	typeMap := map[string][]string{
		"INTEGER":          {"INT", "INT4", "INTEGER"},
		"BIGINT":           {"INT8", "BIGINT", "SERIAL", "SERIAL4", "SERIAL8", "BIGSERIAL"},
		"SERIAL":           {"INT8", "BIGINT", "SERIAL", "SERIAL4", "SERIAL8", "BIGSERIAL"},
		"BIGSERIAL":        {"INT8", "BIGINT", "SERIAL", "SERIAL4", "SERIAL8", "BIGSERIAL"},
		"VARCHAR":          {"VARCHAR", "CHARACTER VARYING"},
		"TEXT":             {"TEXT"},
		"BOOLEAN":          {"BOOL", "BOOLEAN"},
		"TIMESTAMPTZ":      {"TIMESTAMPTZ", "TIMESTAMP WITH TIME ZONE"},
		"DOUBLE PRECISION": {"FLOAT8", "DOUBLE PRECISION"},
	}

	expected = strings.ToUpper(expected)
	actual = strings.ToUpper(actual)

	// Direct match
	if expected == actual {
		return true
	}

	// Check type map
	for _, aliases := range typeMap {
		expectedMatch := false
		actualMatch := false
		for _, alias := range aliases {
			if strings.HasPrefix(expected, alias) {
				expectedMatch = true
			}
			if strings.HasPrefix(actual, alias) {
				actualMatch = true
			}
		}
		if expectedMatch && actualMatch {
			return true
		}
	}

	return false
}
