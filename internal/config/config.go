package config

import (
	"flag"
	"io"
	"os"

	"github.com/loop-kar/pixie/config"
)

const (
	configFileKey       = "configFile"
	defaultFile         = ""
	configFileDescUsage = "Config to be loaded"
)

type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	SMTP     SMTPConfig     `mapstructure:"smtp"`
	Site     SiteConfig     `mapstructure:"site"`
	Config   Config         `mapstructure:"config"`
	Logger   LogConfig      `mapstructure:"logger"`
	S3Config S3Config       `mapstructure:"s3Config"`
}

type ServerConfig struct {
	AppName          string `mapstructure:"appName"`
	LogLevel         string `mapstructure:"logLevel"`
	Port             int    `mapstructure:"port"`
	Host             string `mapstructure:"host"`
	JwtSecretKey     string `mapstructure:"jwtSecretKey"`
	JwtExpiryMinutes int64  `mapstructure:"jwtExpiryMinutes"`
	SecretKey        string `mapstructure:"secretKey"`
	Environment      string `mapstructure:"environment"`
}

type SMTPConfig struct {
	UserName   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Override   bool   `mapstructure:"override"`
	OverrideTo string `mapstructure:"overrideTo"`
}

type SiteConfig struct {
	URLScheme string `mapstructure:"urlScheme"`
	BaseURL   string `mapstructure:"baseUrl"`
}

type Config struct {
	UseJobService bool `mapstructure:"useJobService"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	DBName   string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Schema   string `mapstructure:"schema"`
}

type LogConfig struct {
	License  string `mapstructure:"license"`
	Provider string `mapstructure:"provider"`
	LogLevel string `mapstructure:"logLevel"`
}

type S3Config struct {
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	AccessKeyID     string `mapstructure:"accessKey"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	Endpoint        string `mapstructure:"endpoint"`
	UsePathStyle    bool   `mapstructure:"usePathStyle"`
	ForceHTTPS      bool   `mapstructure:"forceHTTPS"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, configFileKey, defaultFile, configFileDescUsage)

}

func ProvideAppConfig() (AppConfig, error) {

	var configReader io.ReadCloser

	configReader, err := os.Open(configFile)
	if err != nil {
		return AppConfig{}, err
	}

	defer configReader.Close()

	return LoadConfig(configReader)
}

func LoadConfig(configReader io.ReadCloser) (AppConfig, error) {

	cfg := AppConfig{}
	keysToEnvVars := map[string]string{
		"server.port":             "APP_PORT",
		"server.host":             "APP_HOST",
		"server.logLevel":         "LOG_LEVEL",
		"server.appName":          "APP_NAME",
		"server.jwtSecretKey":     "JWT_SECRET_KEY",
		"server.jwtExpiryMinutes": "JWT_EXPIRY_MINUTES",
		"server.secretKey":        "SECRET_KEY",
		"server.AppName":          "APP_NAME",
		"server.environment":      "ENVIRONMENT",

		"database.host":     "DB_HOST",
		"database.name":     "DB_NAME",
		"database.port":     "DB_PORT",
		"database.user":     "DB_USER",
		"database.password": "DB_PASSWORD",
		"database.schema":   "DB_SCHEMA",

		"smtp.username":   "SMTP_USER",
		"smtp.password":   "SMTP_PASSWORD",
		"smtp.host":       "SMTP_HOST",
		"smtp.port":       "SMTP_PORT",
		"smtp.override":   "SMTP_OVERRIDE",
		"smtp.overrideTo": "SMTP_OVERRIDE_TO",

		"site.urlScheme": "SITE_URL_SCHEME",
		"site.baseUrl":   "SITE_BASE_URL",

		"config.useJobService": "USE_JOB_SERVICE",

		"logger.license":  "NEW_RELIC_LICENSE_KEY",
		"logger.provider": "LOG_PROVIDER",
		"logger.logLevel": "LOG_LEVEL",

		"s3Config.region":          "S3_REGION",
		"s3Config.bucket":          "S3_BUCKET",
		"s3Config.accessKey":       "S3_ACCESS_KEY_ID",
		"s3Config.secretAccessKey": "S3_SECRET_ACCESS_KEY",
		"s3Config.endpoint":        "S3_ENDPOINT",
		"s3Config.usePathStyle":    "S3_USE_PATH_STYLE", // true
		"s3Config.forceHTTPS":      "S3_FORCE_HTTPS",    // true
	}

	err := config.LoadConfig(configReader, keysToEnvVars, &cfg)
	cfg.S3Config.Region = "us-east-1"
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
