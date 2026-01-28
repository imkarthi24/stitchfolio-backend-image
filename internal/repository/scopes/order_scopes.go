package scopes

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/loop-kar/pixie/util"
	"gorm.io/gorm"
)

func GetOrders_Search(search string) func(db *gorm.DB) *gorm.DB {
	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&search) {
		return defaultReturn
	}

	return func(db *gorm.DB) *gorm.DB {
		formattedSearch := util.EncloseWithPercentageOperator(search)
		whereClause := fmt.Sprintf(
			`(
				EXISTS (SELECT 1 FROM "stich"."Customers" C WHERE C.id = "stich"."Orders".customer_id AND (C.first_name ILIKE %s OR C.last_name ILIKE %s)) OR 
			 	EXISTS (SELECT 1 FROM "stich"."Users" U WHERE U.id = "stich"."Orders".order_taken_by_id AND (U.first_name ILIKE %s OR U.last_name ILIKE %s)) OR
				"stich"."Orders".id::text ILIKE %s OR
				C.phone_number ILIKE %s OR
			 )`,
			formattedSearch, formattedSearch, formattedSearch, formattedSearch, formattedSearch, formattedSearch,
		)
		return db.Where(whereClause)
	}
}

func GetOrders_Filter(filters string) func(db *gorm.DB) *gorm.DB {

	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&filters) {
		return defaultReturn
	}

	fieldMap := map[string]string{
		"Status":               "status",
		"ExpectedDeliveryDate": "expected_delivery_date",
		"DeliveredDate":        "delivered_date",
		"CustomerId":           "customer_id",
		"OrderTakenById":       "order_taken_by_id",
	}

	return func(db *gorm.DB) *gorm.DB {
		// Split filters - use semicolon when "in" operator has comma-separated values
		// Format examples:
		// - "CustomerId in 1,2,3; OrderTakenById in 5,6"
		// - "Status eq CONFIRMED, ExpectedDeliveryDate btwn '2024-01-01' AND '2024-12-31'"
		var filterArray []string
		if strings.Contains(filters, ";") {
			filterArray = strings.Split(filters, ";")
		} else {
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
				regularFilters = append(regularFilters, filter)
				continue
			}

			// Handle "in" operator for multiple IDs (CustomerId, OrderTakenById)
			if operator == "in" {
				if fieldCamel == "CustomerId" || fieldCamel == "OrderTakenById" {
					values := strings.Split(valueStr, ",")
					var cleanValues []interface{}
					for _, v := range values {
						v = strings.TrimSpace(v)
						if v != "" {
							cleanValues = append(cleanValues, v)
						}
					}
					if len(cleanValues) > 0 {
						inFilters[dbField] = cleanValues
					}
				} else {
					regularFilters = append(regularFilters, filter)
				}
			} else {
				// Regular filter (eq, lt, gt, neq, btwn, etc.)
				regularFilters = append(regularFilters, filter)
			}
		}

		if len(regularFilters) > 0 {
			regularFiltersStr := strings.Join(regularFilters, ",")
			queryString := util.BuildQuery(regularFiltersStr, fieldMap)
			if !util.IsNilOrEmptyString(&queryString) {
				db = db.Where(queryString)
			}
		}

		for dbField, values := range inFilters {
			if len(values) > 0 {
				db = db.Where(fmt.Sprintf(`"%s" IN ?`, dbField), values)
			}
		}

		return db
	}
}
