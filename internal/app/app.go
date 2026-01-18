package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	config_cache "github.com/imkarthi24/sf-backend/internal/cache"
	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/pkg/db/migrator"
	pkgLog "github.com/imkarthi24/sf-backend/pkg/log"
	"github.com/newrelic/go-agent/v3/newrelic"

	"gorm.io/gorm"
)

type App struct {
	Server                 *gin.Engine
	AppConfig              config.AppConfig
	StitchDB               *gorm.DB
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
	dbConn, err := a.StitchDB.DB()
	checkErr(err)

	err = dbConn.Close()
	checkErr(err)

	//	pkgLog.FromCtx(ctx).Info("App Service Stopped...")
}

func (a *App) Migrate(ctx *context.Context, checkErr func(err error)) {
	migrator := migrator.NewMigrator(a.StitchDB)

	entityList := []interface{}{
		// &entities.Channel{},
		// &entities.Customer{},
		// &entities.DressType{},
		// &entities.EmailNotification{},
		// &entities.EnquiryHistory{},
		// &entities.Enquiry{},
		// &entities.MasterConfig{},
		// &entities.Measurement{},
		// &entities.MeasurementHistory{},
		// &entities.Notification{},
		// &entities.OrderHistory{},
		// &entities.Order{},
		// &entities.OrderItem{},
		// &entities.Person{},
		// &entities.UserChannelDetail{},
		// &entities.UserConfig{},
		// &entities.User{},
		// &entities.WhatsappNotification{},
	}

	//************************//
	/*
	 0. Use Migrate only for initial migration
	 1. Make changes to entity files
	 2. Provide proper migration name in GenerateAlterMigration with incremented number and snake case eg:" 002_person_entity_update"
	 4. Only uncomment entity which is changed in entityList to generate minimal alter migration file
	 5. Comment back the GenerateAlterMigration line after generating migration file to avoid overwriting
	 6. Use GenerateAlterMigration to generate alter migration file
	*/
	//************************//

	//migrator.Migrate(entityList, checkErr)

	migrator.GenerateAlterMigration(entityList, "002_person_entity_update")
}
