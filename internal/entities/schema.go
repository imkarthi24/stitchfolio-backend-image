package entities

import "os"

var dbSchema string

// GetSchema returns the database schema name
// It first checks if it's been set via InitSchema, otherwise reads from DB_SCHEMA env var
func GetSchema() string {
	if dbSchema != "" {
		return dbSchema
	}
	// Fallback to environment variable
	if schema := os.Getenv("DB_SCHEMA"); schema != "" {
		return schema
	}
	// Default fallback (should not be used in production)
	return "public"
}

// InitSchema initializes the schema name from config
// This should be called during application startup
func InitSchema(schema string) {
	dbSchema = schema
}

// TableNameWithSchema returns a table name with schema prefix
func TableNameWithSchema(tableName string) string {
	return GetSchema() + "." + tableName
}

// TableNameForQueryWithSchema returns a table name formatted for queries with schema prefix
func TableNameForQueryWithSchema(tableName string) string {
	return "\"" + GetSchema() + "\".\"" + tableName + "\" E"
}
