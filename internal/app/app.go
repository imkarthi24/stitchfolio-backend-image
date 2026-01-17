package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	config_cache "github.com/imkarthi24/sf-backend/internal/cache"
	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository"
	pkgLog "github.com/imkarthi24/sf-backend/pkg/log"
	"github.com/newrelic/go-agent/v3/newrelic"

	"gorm.io/gorm"
)

type App struct {
	Server                 *gin.Engine
	AppConfig              config.AppConfig
	ChitDb                 *gorm.DB
	NewRelic               *newrelic.Application
	MasterConfigRepository repository.MasterConfigRepository
}

func (a *App) Start(ctx *context.Context, checkErr func(err error)) {
	go func(ctx *context.Context) {

		//App startup essentials
		// entities.InitSchema(a.AppConfig.Database.Schema)
		pkgLog.InitLogger(a.NewRelic)
		config_cache.InitMasterConfig(a.MasterConfigRepository)

		//Spin up server
		//Render needs host to be empty , while Oracel needs it as 0.0.0.0. so trim by setting " " for render

		host := a.AppConfig.Server.Host
		if strings.Contains(host, "render") {
			host = ""
		}
		address := fmt.Sprintf("%s:%d", host, a.AppConfig.Server.Port)
		err := a.Server.Run(address)
		if err != nil {
			checkErr(err)
		}
	}(ctx)
}

func (a *App) Shutdown(ctx *context.Context, checkErr func(err error)) {
	dbConn, err := a.ChitDb.DB()
	checkErr(err)

	err = dbConn.Close()
	checkErr(err)

	//	pkgLog.FromCtx(ctx).Info("App Service Stopped...")
}

func (a *App) Migrate(ctx *context.Context, checkErr func(err error)) {

	entityList := []interface{}{
		// &entities.User{},
		// &entities.Channel{},
		// &entities.Notification{},
		// &entities.MasterConfig{},
		// &entities.EmailNotification{},
		// &entities.WhatsappNotification{},
		// &entities.Enquiry{},
		// &entities.EnquiryHistory{},
		// &entities.Customer{},
		&entities.Order{},
		// &entities.OrderItem{},
		// &entities.Measurement{},
	}
	for _, entity := range entityList {
		// pass entity directly (it's already a pointer to a struct)
		err := a.ChitDb.AutoMigrate(entity)
		if err != nil {
			fmt.Printf("Migration error for %T: %v\n", entity, err)
			checkErr(err)
		}
	}

}
