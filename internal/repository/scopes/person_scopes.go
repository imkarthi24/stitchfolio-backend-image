package scopes

import (
	"fmt"

	"github.com/imkarthi24/sf-backend/pkg/util"
	"gorm.io/gorm"
)

func PersonNameFilter(name string) func(db *gorm.DB) *gorm.DB {

	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&name) {
		return defaultReturn
	}

	return func(db *gorm.DB) *gorm.DB {
		enclosedName := util.EncloseWithPercentageOperator(name)
		whereClause := fmt.Sprintf("(name ILIKE  %s)", enclosedName)
		return db.Where(whereClause)
	}
}
