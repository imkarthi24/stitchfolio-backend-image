package scopes

import (
	"fmt"
	"strings"

	"github.com/imkarthi24/sf-backend/pkg/util"
	"gorm.io/gorm"
)

func GetOrders_Filter(filters string) func(db *gorm.DB) *gorm.DB {

	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&filters) {
		return defaultReturn
	}

	fieldMap := map[string]string{
		"Status":               "status",
		"ExpectedDeliveryDate": "expected_delivery_date",
		"DeliveredDate":        "delivered_date",
	}

	return func(db *gorm.DB) *gorm.DB {
		// Parse filters to separate DressType from other filters
		filterArray := strings.Split(filters, ",")
		var orderFilters []string
		var dressTypeFilter string

		for _, filter := range filterArray {
			parts := strings.Split(strings.TrimSpace(filter), " ")
			if len(parts) >= 3 {
				field := strings.ToLower(parts[0])
				if field == "name" {
					dressTypeFilter = filter
				} else {
					orderFilters = append(orderFilters, filter)
				}
			}
		}

		orderFiltersStr := strings.Join(orderFilters, ",")
		if !util.IsNilOrEmptyString(&orderFiltersStr) {
			queryString := util.BuildQuery(orderFiltersStr, fieldMap)
			if !util.IsNilOrEmptyString(&queryString) {
				db = db.Where(queryString)
			}
		}

		// Handle DressType filter separately (requires JOIN through OrderItems -> Measurements)
		if !util.IsNilOrEmptyString(&dressTypeFilter) {
			parts := strings.Split(strings.TrimSpace(dressTypeFilter), " ")
			if len(parts) >= 3 {
				operator := parts[1]
				value := strings.Join(parts[2:], " ")

				var operatorSymbol string
				switch operator {
				case "eq":
					operatorSymbol = "="
				case "lt":
					operatorSymbol = "<"
				case "gt":
					operatorSymbol = ">"
				case "neq":
					operatorSymbol = "!="
				default:
					operatorSymbol = "="
				}

				existsQuery := fmt.Sprintf(
					`EXISTS (SELECT 1 FROM "stich"."OrderItems" OI 
				 JOIN "stich"."Measurements" M ON M.id = OI.measurement_id 
				 JOIN "stich"."DressTypes" DT ON DT.id = M.dress_type_id 
				 WHERE OI.order_id = "stich"."Orders".id AND DT.name %s ?)`,
					operatorSymbol,
				)
				//TODO : imporve
				db = db.Where(existsQuery, value)
			}
		}

		return db
	}
}
