package utils

import (
	"context"
	"fmt"

	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/constants"
	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/model/models"
	"github.com/loop-kar/pixie/util"

	pkgConst "github.com/loop-kar/pixie/constants"
)

func GetSiteURL(config config.SiteConfig) string {
	return fmt.Sprintf("%s://%s/", config.URLScheme, config.BaseURL)
}

func GetHealthEndpoint(config config.ServerConfig) string {
	return fmt.Sprintf("%s/%s%s", config.Host, constants.API_PREFIX_V1, constants.HEALTH)
}

func GetChannelId(ctx *context.Context) uint {
	session := GetSession(ctx)
	if session == nil {
		return 0
	}

	return session.ChannelId
}

func GetRole(ctx *context.Context) entities.RoleType {
	session := GetSession(ctx)

	if session == nil {
		return "UNKNOWN"
	}

	return session.Role
}

func GetUserId(ctx *context.Context) uint {
	session := GetSession(ctx)
	if session == nil {
		return 0
	}

	return *session.UserId
}

func GetAccessibleLocationIds(ctx *context.Context) []uint {
	session := GetSession(ctx)
	if session == nil {
		return []uint{}
	}

	return session.AccessibleLocationIds
}

func GetSession(ctx *context.Context) *models.Session {
	val := util.ReadValueFromContext(ctx, pkgConst.SESSION)
	session, ok := val.(*models.Session)
	if !ok {
		return nil
	}

	return session

}
