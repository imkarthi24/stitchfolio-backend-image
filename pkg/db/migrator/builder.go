package migrator

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func generateSQL(db *gorm.DB, models []interface{}) ([]string, error) {
	var sqlStatements []string
	var foreignKeyStatements []string

	for _, model := range models {
		// Parse the model to get schema information
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return nil, err
		}

		// Build CREATE TABLE SQL manually using schema info
		sql := buildCreateTableSQL(stmt.Schema)
		sqlStatements = append(sqlStatements, sql)

		// Build index SQLs
		indexSQLs := buildIndexSQL(stmt.Schema)
		sqlStatements = append(sqlStatements, indexSQLs...)

		// Collect foreign key constraints (to be added after all tables are created)
		fkSQLs := buildForeignKeySQL(stmt.Schema)
		foreignKeyStatements = append(foreignKeyStatements, fkSQLs...)
	}

	// Add foreign keys at the end (after all tables exist)
	sqlStatements = append(sqlStatements, foreignKeyStatements...)

	return sqlStatements, nil
}

func buildCreateTableSQL(schema *schema.Schema) string {
	var columns []string
	var primaryKeys []string
	var tableConstraints []string

	for _, field := range schema.Fields {
		if field.DBName == "" {
			continue
		}
		column := buildColumnDefinition(field)
		columns = append(columns, column)

		if field.PrimaryKey {
			primaryKeys = append(primaryKeys, field.DBName)
		}

		// Add check constraints if defined
		if checkConstraint := getCheckConstraint(field); checkConstraint != "" {
			tableConstraints = append(tableConstraints, checkConstraint)
		}
	}

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s",
		getTableNameWithQuotes(schema),
		strings.Join(columns, ",\n  "))

	if len(primaryKeys) > 0 {
		sql += fmt.Sprintf(",\n  PRIMARY KEY (%s)", strings.Join(primaryKeys, ", "))
	}

	// Add table-level constraints
	for _, constraint := range tableConstraints {
		sql += ",\n  " + constraint
	}

	sql += "\n)"

	return sql
}

func buildColumnDefinition(field *schema.Field) string {
	var parts []string

	// Column name
	parts = append(parts, field.DBName)

	// Data type
	dataType := getPostgresType(field)
	parts = append(parts, dataType)

	// NOT NULL
	if field.NotNull {
		parts = append(parts, "NOT NULL")
	}

	// UNIQUE
	if field.Unique {
		parts = append(parts, "UNIQUE")
	}

	// DEFAULT
	if field.DefaultValue != "" {
		if field.DataType == "text" && !strings.HasPrefix(field.DefaultValue, "'") {
			//Add single quotes for text default values
			field.DefaultValue = fmt.Sprintf("'%s'", field.DefaultValue)
		}
		parts = append(parts, fmt.Sprintf("DEFAULT %s", field.DefaultValue))
	}

	// CHECK constraint (column-level)
	if check, ok := field.TagSettings["CHECK"]; ok {
		parts = append(parts, fmt.Sprintf("CHECK (%s)", check))
	}

	return strings.Join(parts, " ")
}

// getTableNameWithQuotes returns the table name with double quotes if applicable
// if `schema.Table` is `stitch.Channels`, it returns `stitch."Channels"`
func getTableNameWithQuotes(schema *schema.Schema) string {
	schemaParts := strings.Split(schema.Table, ".")
	if len(schemaParts) > 1 {
		return fmt.Sprintf(`%s."%s"`, schemaParts[0], schemaParts[1])
	}
	return schema.Table
}

func getCheckConstraint(field *schema.Field) string {
	// Check for table-level CHECK constraints
	if check, ok := field.TagSettings["CONSTRAINT"]; ok {
		constraintName := fmt.Sprintf("chk_%s_%s", field.Schema.Table, field.DBName)
		return fmt.Sprintf("CONSTRAINT %s CHECK (%s)", constraintName, check)
	}
	return ""
}

func getPostgresType(field *schema.Field) string {
	// Check for custom size
	if size, ok := field.TagSettings["SIZE"]; ok {
		if field.DataType == "string" {
			return fmt.Sprintf("VARCHAR(%s)", size)
		}
	}

	// Check for custom type
	if fieldType, ok := field.TagSettings["TYPE"]; ok {
		return strings.ToUpper(fieldType)
	}

	// Auto increment
	if field.AutoIncrement {
		switch field.Size {
		case 16:
			return "SMALLSERIAL"
		case 64:
			return "BIGSERIAL"
		default:
			return "SERIAL"
		}
	}

	// Default type mapping
	switch field.DataType {
	case "bool":
		return "BOOLEAN"
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32":
		return "INTEGER"
	case "int64", "uint64":
		return "BIGINT"
	case "float32":
		return "REAL"
	case "float64":
		return "DOUBLE PRECISION"
	case "string":
		return "TEXT"
	case "time.Time":
		return "TIMESTAMPTZ"
	case "*time.Time":
		return "TIMESTAMPTZ"
	case "[]byte":
		return "BYTEA"
	default:
		return "TEXT"
	}
}

func buildIndexSQL(schema *schema.Schema) []string {
	var sqls []string

	for _, index := range schema.ParseIndexes() {
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
			getTableNameWithQuotes(schema),
			strings.Join(fields, ", "))

		sqls = append(sqls, sql)
	}

	return sqls
}

func buildForeignKeySQL(schema *schema.Schema) []string {
	var sqls []string

	for _, field := range schema.Fields {
		// Check if field has a foreign key relation
		for _, rel := range schema.Relationships.Relations {
			if rel.Type == "belongs_to" && len(rel.References) > 0 {
				// Find if this field is part of the foreign key
				for _, ref := range rel.References {
					if ref.ForeignKey.DBName == field.DBName {
						fkSQL := buildSingleForeignKey(schema, rel, ref)
						if fkSQL != "" {
							sqls = append(sqls, fkSQL)
						}
					}
				}
			}
		}
	}

	return sqls
}

func buildSingleForeignKey(schema *schema.Schema, rel *schema.Relationship, ref *schema.Reference) string {
	// Get constraint settings from field tags
	onDelete := "RESTRICT"
	onUpdate := "RESTRICT"

	// Check for constraint tag settings
	if rel.Field != nil {
		if constraintTag, ok := rel.Field.TagSettings["CONSTRAINT"]; ok {
			parts := strings.Split(constraintTag, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "OnDelete:") {
					onDelete = strings.TrimPrefix(part, "OnDelete:")
				} else if strings.HasPrefix(part, "OnUpdate:") {
					onUpdate = strings.TrimPrefix(part, "OnUpdate:")
				}
			}
		}
	}

	// Generate constraint name
	fkName := fmt.Sprintf("fk_%s_%s", schema.Name, ref.ForeignKey.DBName)

	// Build foreign key SQL
	sql := fmt.Sprintf(
		"ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE %s ON UPDATE %s",
		getTableNameWithQuotes(schema),
		fkName,
		ref.ForeignKey.DBName,
		getTableNameWithQuotes(rel.FieldSchema),
		ref.PrimaryKey.DBName,
		onDelete,
		onUpdate,
	)

	return sql
}
