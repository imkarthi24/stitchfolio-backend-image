package scopes

import (
	"fmt"
	"strings"

	"github.com/loop-kar/pixie/constants"
	"github.com/loop-kar/pixie/util"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

const CHANNEL_ID = "channel_id"
const And = " AND "
const OR = " OR "

func IsActive(params ...string) func(db *gorm.DB) *gorm.DB {

	if len(params) == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active", true)
		}
	}

	//chain the where conditions and return

	var IsActive string = "is_active"
	var And string = " AND "

	return func(db *gorm.DB) *gorm.DB {

		whereClause := ""
		funk.ForEach(params, func(param string) {
			if !strings.HasPrefix(param, "E") {
				param = util.EncloseWithSymbol(param, "\"")
			}

			whereClause = whereClause + fmt.Sprintf("%s.%s = true", param, IsActive) + And
		})

		//removing the last 'and' word
		whereClause = strings.TrimSuffix(whereClause, And)

		return db.Where(whereClause)
	}

}

func Channel(params ...string) func(db *gorm.DB) *gorm.DB {

	if len(params) == 0 {
		return func(db *gorm.DB) *gorm.DB {

			var channelId uint
			if id, ok := db.Get(constants.CHANNEL_ID); ok {
				channelId = id.(uint)
			}

			//System Admin needs access to all Data
			if channelId == 0 {
				return db
			}

			return db.Where("channel_id", channelId)

		}
	}

	//chain the where conditions and return

	return func(db *gorm.DB) *gorm.DB {

		var channelId uint
		if id, ok := db.Get(constants.CHANNEL_ID); ok {
			channelId = id.(uint)
		}

		//System Admin needs access to all Data
		if channelId == 0 {
			return db
		}

		whereClause := ""
		funk.ForEach(params, func(param string) {
			if !strings.HasPrefix(param, "E") {
				param = util.EncloseWithSymbol(param, "\"")
			}

			whereClause = whereClause + fmt.Sprintf("%s.%s = %d", param, CHANNEL_ID, channelId) + And
		})

		//removing the last 'and' word
		whereClause = strings.Trim(whereClause, And)

		return db.Where(whereClause)
	}

}

func AccessibleChannels(accesibleChannelIds []uint) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		var channelId uint
		if id, ok := db.Get(constants.CHANNEL_ID); ok {
			channelId = id.(uint)
		}

		//System Admin needs access to all Data
		if channelId == 0 {
			return db
		}

		return db.Where("channel_id", accesibleChannelIds)

	}

}

func ILike(query string, params ...string) func(db *gorm.DB) *gorm.DB {

	if len(params) == 0 || util.IsNilOrEmptyString(&query) {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	return func(db *gorm.DB) *gorm.DB {

		whereClause := ""
		query = util.EncloseWithPercentageOperator(query)
		funk.ForEach(params, func(param string) {

			whereClause = whereClause + fmt.Sprintf(`%s ILIKE %s`, param, query) + OR
		})

		//removing the last 'and' word
		whereClause = strings.Trim(whereClause, OR)

		return db.Where(whereClause)
	}

}

func SelectFields(params ...string) func(db *gorm.DB) *gorm.DB {

	if len(params) == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Select("id", params)
	}

}
