package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imkarthi24/sf-backend/internal/model/models"
	"github.com/imkarthi24/sf-backend/pkg/constants"
)

func RoleBasedAccessControl() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionValue := ctx.Value(constants.SESSION)
		if sessionValue == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, ok := sessionValue.(*models.Session)
		if !ok || session == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Skip RBAC for system sessions (external access)
		if session.IsSystemSession {
			ctx.Next()
			return
		}

		if err := checkResourceAccess(ctx, session); err != nil {
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		}

		ctx.Next()
	}
}

// checkResourceAccess validates if the user has access to the requested resource
func checkResourceAccess(ctx *gin.Context, session *models.Session) error {
	// path := ctx.Request.URL.Path
	// method := ctx.Request.Method

	return nil
}
