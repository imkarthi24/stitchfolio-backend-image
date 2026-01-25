package scopes

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/loop-kar/pixie/util"
	"gorm.io/gorm"
)

func GetExpenseTrackers_Search(search string) func(db *gorm.DB) *gorm.DB {
	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&search) {
		return defaultReturn
	}

	return func(db *gorm.DB) *gorm.DB {
		formattedSearch := util.EncloseWithPercentageOperator(search)
		whereClause := fmt.Sprintf(
			`("stich"."ExpenseTrackers".bill_number ILIKE %s OR 
			 "stich"."ExpenseTrackers".company_name ILIKE %s OR 
			 "stich"."ExpenseTrackers".material ILIKE %s OR 
			 "stich"."ExpenseTrackers".notes ILIKE %s)`,
			formattedSearch, formattedSearch, formattedSearch, formattedSearch,
		)
		return db.Where(whereClause)
	}
}

func GetExpenseTrackers_Filter(filters string) func(db *gorm.DB) *gorm.DB {
	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&filters) {
		return defaultReturn
	}

	fieldMap := map[string]string{
		"PurchaseDate": "purchase_date",
		"Price":        "price",
		"Location":     "location",
	}

	return func(db *gorm.DB) *gorm.DB {
		// Split filters - use semicolon when "in" operator has comma-separated values
		// Format examples:
		// - "PurchaseDate in '2024-01-01','2024-01-02','2024-01-03'"
		// - "PurchaseDate btwn '2024-01-01' AND '2024-12-31'"
		var filterArray []string
		if strings.Contains(filters, ";") {
			filterArray = strings.Split(filters, ";")
		} else {
			// For backward compatibility, try comma separation
			filterArray = strings.Split(filters, ",")
		}

		var regularFilters []string
		inFilters := make(map[string][]interface{})

		for _, filter := range filterArray {
			filter = strings.TrimSpace(filter)
			if filter == "" {
				continue
			}

			parts := strings.Fields(filter)
			if len(parts) < 3 {
				continue
			}

			field := parts[0]
			operator := strings.ToLower(parts[1])
			valueStr := strings.Join(parts[2:], " ")

			// Normalize field name to CamelCase for lookup
			fieldCamel := strcase.ToCamel(field)
			dbField, exists := fieldMap[fieldCamel]
			if !exists {
				// Unknown field, treat as regular filter
				regularFilters = append(regularFilters, filter)
				continue
			}

			// Handle "in" operator for multiple PurchaseDate values
			if operator == "in" {
				// Only allow "in" for PurchaseDate
				if fieldCamel == "PurchaseDate" {
					// Parse comma-separated date values
					values := strings.Split(valueStr, ",")
					var cleanValues []interface{}
					for _, v := range values {
						v = strings.TrimSpace(v)
						// Remove quotes if present
						v = strings.Trim(v, "'\"")
						if v != "" {
							cleanValues = append(cleanValues, v)
						}
					}
					if len(cleanValues) > 0 {
						inFilters[dbField] = cleanValues
					}
				} else {
					// "in" not supported for other fields, treat as regular filter
					regularFilters = append(regularFilters, filter)
				}
			} else {
				// Regular filter (eq, lt, gt, neq, btwn, etc.)
				regularFilters = append(regularFilters, filter)
			}
		}

		// Apply regular filters using BuildQuery (handles PurchaseDate btwn, Price, etc.)
		if len(regularFilters) > 0 {
			regularFiltersStr := strings.Join(regularFilters, ",")
			queryString := util.BuildQuery(regularFiltersStr, fieldMap)
			if !util.IsNilOrEmptyString(&queryString) {
				db = db.Where(queryString)
			}
		}

		// Apply IN filters for multiple PurchaseDate values
		for dbField, values := range inFilters {
			if len(values) > 0 {
				// Use GORM's Where with IN clause for dates
				// Format: purchase_date IN ('2024-01-01', '2024-01-02', '2024-01-03')
				db = db.Where(fmt.Sprintf(`"%s" IN ?`, dbField), values)
			}
		}

		return db
	}
}
