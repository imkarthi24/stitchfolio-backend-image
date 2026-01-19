package scopes

import (
	"fmt"

	"github.com/loop-kar/pixie/util"
	"gorm.io/gorm"
)

func SearchNameOrEmailOrPhone_Filter(name string) func(db *gorm.DB) *gorm.DB {

	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&name) {
		return defaultReturn
	}

	return func(db *gorm.DB) *gorm.DB {
		enclosedName := util.EncloseWithPercentageOperator(name)
		whereClause := fmt.Sprintf("(first_name ILIKE %s OR last_name ILIKE  %s OR email ILIKE %s OR phone_number ILIKE %s )", enclosedName, enclosedName, enclosedName, enclosedName)
		return db.Where(whereClause)
	}
}

// status eq xyz , or date eq , or name
func GetEnquiries_Filter(filters string) func(db *gorm.DB) *gorm.DB {

	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&filters) {
		return defaultReturn
	}

	fieldMap := map[string]string{
		"Source":              "source",
		"ReferredBy":          "referred_by",
		"ReferrerPhoneNumber": "referrer_phone_number",
		"Course":              "course",
		"PlacementRequired":   "placement_required",
		"CurrentlyWorking":    "currenlty_working",
		"Profession":          "profession",
		"Status":              "status",
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where(util.BuildQuery(filters, fieldMap))
	}
}
