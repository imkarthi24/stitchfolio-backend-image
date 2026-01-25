package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/model/models"
	"github.com/loop-kar/pixie/constants"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"
	"github.com/mitchellh/mapstructure"
)

func VerifyJWT(secretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if ctx.Request.Header["Token"] == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := util.VerifyJWT(ctx.Request.Header["Token"][0], secretKey)

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		//TODO: can this be passed as an interface and then decoded/encoded?
		sessionDetails := &models.Session{}
		mapstructure.Decode(token.Claims, sessionDetails)
		ctx.Set(constants.SESSION, sessionDetails)

		err = restrictDevAccess(ctx, *sessionDetails)
		if err != nil {
			ctx.AbortWithError(http.StatusForbidden, err)
		}
		ctx.Next()
	}
}

func GenerateExternalSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		channelHeaders, channelExists := ctx.Request.Header["Channel"]
		channelIdHeaders, channelIdExists := ctx.Request.Header["Channel-Id"]

		if !channelExists || len(channelHeaders) == 0 || !channelIdExists || len(channelIdHeaders) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errs.NewXError(errs.INVALID_REQUEST, "missing required headers", nil))
			return
		}

		channel := channelHeaders[0]
		channelId, err := strconv.Atoi(channelIdHeaders[0])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errs.NewXError(errs.INVALID_REQUEST, "invalid Channel-Id header", nil))
			return
		}

		if util.IsNilOrEmptyString(&channel) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errs.NewXError(errs.INVALID_REQUEST, "invalid channel header", nil))
		}
		sessionDetails := &models.Session{
			Role:            entities.SUPERADMIN, //TODO: outsourced to freelance tailors
			ChannelName:     channel,
			ChannelId:       uint(channelId),
			IsSystemSession: true,
		}
		ctx.Set(constants.SESSION, sessionDetails)
		ctx.Next()
	}
}

func restrictDevAccess(ctx *gin.Context, session models.Session) error {

	//CHECK DEV role
	if session.Role == entities.DEV {
		if ctx.Request.Method != "GET" {
			return errors.New("action not allowed for dev accounts")
		}
	}

	return nil
}
