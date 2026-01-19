package di

import (
	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/loop-kar/pixie/db"
	pkgservice "github.com/loop-kar/pixie/service"
	pkgemail "github.com/loop-kar/pixie/service/email"
)

// ProvideServiceContainer builds the shared "pkg/service" dependency container.
// This is intended to be injected into internal services via Wire.
func ProvideServiceContainer(appConfig config.AppConfig) *pkgservice.Service {
	emailSvc := pkgemail.NewEmailService(pkgemail.SMTPConfig{
		UserName:   appConfig.SMTP.UserName,
		Password:   appConfig.SMTP.Password,
		Host:       appConfig.SMTP.Host,
		Port:       appConfig.SMTP.Port,
		Override:   appConfig.SMTP.Override,
		OverrideTo: appConfig.SMTP.OverrideTo,
	})

	return pkgservice.NewService(
		pkgservice.WithEmailService(emailSvc),
	)
}

// ProvideDatabaseConnectionParams maps the internal config to the database connection params
func ProvideDatabaseConnectionParams(dbConfig config.DatabaseConfig) db.DatabaseConnectionParams {
	sslMode := "prefer"

	return db.DatabaseConnectionParams{
		Host:     dbConfig.Host,
		Port:     dbConfig.Port,
		Username: dbConfig.Username,
		DBName:   dbConfig.DBName,
		Password: dbConfig.Password,
		SSLMode:  sslMode,
		Schema:   dbConfig.Schema,
	}
}
