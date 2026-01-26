package scopes

import (
	"fmt"

	"github.com/loop-kar/pixie/util"
	"gorm.io/gorm"
)

// GetMeasurements_Search - OLD SEARCH (COMMENTED OUT)
// This search searched by DressType name, Person name, and TakenBy name
// func GetMeasurements_Search(search string) func(db *gorm.DB) *gorm.DB {
// 	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }
//
// 	if util.IsNilOrEmptyString(&search) {
// 		return defaultReturn
// 	}
//
// 	return func(db *gorm.DB) *gorm.DB {
// 		formattedSearch := util.EncloseWithPercentageOperator(search)
// 		whereClause := fmt.Sprintf(
// 			`(EXISTS (SELECT 1 FROM "stich"."DressTypes" DT WHERE DT.id = "stich"."Measurements".dress_type_id AND DT.name ILIKE %s) OR
// 			 EXISTS (SELECT 1 FROM "stich"."Persons" P WHERE P.id = "stich"."Measurements".person_id AND (P.first_name ILIKE %s OR P.last_name ILIKE %s))
// 			 -- OR EXISTS (SELECT 1 FROM "stich"."Users" U WHERE U.id = "stich"."Measurements".taken_by_id AND (U.first_name ILIKE %s OR U.last_name ILIKE %s))
// 			 )`,
// 			formattedSearch, formattedSearch, formattedSearch, formattedSearch, formattedSearch,
// 		)
// 		return db.Where(whereClause)
// 	}
// }

// GetMeasurements_SearchByCustomerName - NEW SEARCH
// Searches by Customer Name and returns all Persons of that customer
func GetMeasurements_SearchByCustomerName(search string) func(db *gorm.DB) *gorm.DB {
	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&search) {
		return defaultReturn
	}

	return func(db *gorm.DB) *gorm.DB {
		formattedSearch := util.EncloseWithPercentageOperator(search)
		whereClause := fmt.Sprintf(
			`EXISTS (SELECT 1 FROM "stich"."Persons" P 
				INNER JOIN "stich"."Customers" C ON C.id = P.customer_id 
				WHERE P.id = "stich"."Measurements".person_id 
				AND (C.first_name ILIKE %s OR C.last_name ILIKE %s OR CONCAT(C.first_name, ' ', C.last_name) ILIKE %s))`,
			formattedSearch, formattedSearch, formattedSearch,
		)
		return db.Where(whereClause)
	}
}

// GetMeasurements_Filter - OLD FILTER (COMMENTED OUT)
// This filter supported PersonId, TakenById, and DressTypeId with IN operator
// func GetMeasurements_Filter(filters string) func(db *gorm.DB) *gorm.DB {
// 	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }
//
// 	if util.IsNilOrEmptyString(&filters) {
// 		return defaultReturn
// 	}
//
// 	fieldMap := map[string]string{
// 		"PersonId":    "person_id",
// 		"TakenById":   "taken_by_id",
// 		"DressTypeId": "dress_type_id",
// 	}
//
// 	return func(db *gorm.DB) *gorm.DB {
// 		// Split filters - use semicolon when "in" operator has comma-separated values
// 		// Format examples:
// 		// - "PersonId in 1,2,3; TakenById in 5,6"
// 		var filterArray []string
// 		if strings.Contains(filters, ";") {
// 			filterArray = strings.Split(filters, ";")
// 		}
//
// 		var regularFilters []string
// 		inFilters := make(map[string][]interface{})
//
// 		for _, filter := range filterArray {
// 			filter = strings.TrimSpace(filter)
// 			if filter == "" {
// 				continue
// 			}
//
// 			parts := strings.Fields(filter)
// 			if len(parts) < 3 {
// 				continue
// 			}
//
// 			field := parts[0]
// 			operator := strings.ToLower(parts[1])
// 			valueStr := strings.Join(parts[2:], " ")
//
// 			// Normalize field name to CamelCase
// 			dbField, exists := fieldMap[field]
// 			if !exists {
// 				regularFilters = append(regularFilters, filter)
// 				continue
// 			}
//
// 			// Handle "in" operator for multiple IDs
// 			if operator == "in" {
// 				values := strings.Split(valueStr, ",")
// 				var cleanValues []interface{}
// 				for _, v := range values {
// 					v = strings.TrimSpace(v)
// 					if v != "" {
// 						cleanValues = append(cleanValues, v)
// 					}
// 				}
// 				if len(cleanValues) > 0 {
// 					inFilters[dbField] = cleanValues
// 				}
// 			} else {
// 				regularFilters = append(regularFilters, filter)
// 			}
// 		}
//
// 		// Apply regular filters using existing BuildQuery
// 		if len(regularFilters) > 0 {
// 			regularFiltersStr := strings.Join(regularFilters, ",")
// 			queryString := util.BuildQuery(regularFiltersStr, fieldMap)
// 			if !util.IsNilOrEmptyString(&queryString) {
// 				db = db.Where(queryString)
// 			}
// 		}
//
// 		// Apply IN filters using GORM's Where with IN clause
// 		for dbField, values := range inFilters {
// 			if len(values) > 0 {
// 				// Use GORM's Where with IN clause
// 				db = db.Where(fmt.Sprintf(`"%s" IN ?`, dbField), values)
// 			}
// 		}
//
// 		return db
// 	}
// }
