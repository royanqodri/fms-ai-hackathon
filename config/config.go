package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var (
	config *Config
)

// option defines configuration option
type option struct {
}

// Init initializes `config` from the default config file.
// use `WithConfigFile` to specify the location of the config file
func Init() error {
	// Create viper instance
	v := viper.New()

	// Load environment variables from OS first
	v.AutomaticEnv()

	// Determine mode: local or prod based on ENV var
	mode := v.GetString("SERVICE_ENV")
	if mode == "" {
		mode = "LOCAL" // Default mode if not set
	}

	if mode != "PROD" {
		// Local dev mode using config.env file
		v.SetConfigName("config") // config.env without extension
		v.SetConfigType("env")
		v.AddConfigPath(".")        // current directory
		v.AddConfigPath("./config") // optional if you keep config files in /config

		// Read config file
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	v.SetDefault("SERVICE_NAME", "gateway")
	v.SetDefault("VERSION_APP", "latest")

	// Initialize config struct
	config = &Config{}

	// Service Config
	config.Service = ServiceConfig{
		Name:         v.GetString("SERVICE_NAME"),
		Environment:  v.GetString("SERVICE_ENV"),
		HttpAddress:  v.GetString("SERVICE_HTTP_ADDRESS"),
		LogLevel:     v.GetString("SERVICE_LOG_LEVEL"),
		ReadTimeout:  v.GetInt("SERVICE_READ_TIMEOUT"),
		WriteTimeout: v.GetInt("SERVICE_WRITE_TIMEOUT"),
		IdleTimeout:  v.GetInt("SERVICE_IDLE_TIMEOUT"),
		Location:     v.GetString("SERVICE_LOCATION"),
	}

	// DB Master Config
	config.DBMasterConfig = DBConfig{
		Driver:          v.GetString("DB_MASTER_DRIVER"),
		Host:            v.GetString("DB_MASTER_HOST"),
		Port:            v.GetString("DB_MASTER_PORT"),
		DatabaseName:    v.GetString("DB_MASTER_DATABASE_NAME"),
		User:            v.GetString("DB_MASTER_USER"),
		Password:        v.GetString("DB_MASTER_PASSWORD"),
		SslMode:         v.GetString("DB_MASTER_SSL_MODE"),
		Timezone:        v.GetString("DB_MASTER_TIMEZONE"),
		MaxIdleConn:     v.GetInt("DB_MASTER_MAX_IDLE_CONN"),
		MaxOpenConn:     v.GetInt("DB_MASTER_MAX_OPEN_CONN"),
		MaxConnLifetime: time.Duration(v.GetInt("DB_MASTER_MAX_CONN_LIFE_TIME")) * time.Second,
		MaxConnIdletime: time.Duration(v.GetInt("DB_MASTER_MAX_CONN_IDLE_TIME")) * time.Second,
	}

	// DB Operation Config
	config.DBOperationConfig = DBConfig{
		Driver:          v.GetString("DB_OPERATION_DRIVER"),
		Host:            v.GetString("DB_OPERATION_HOST"),
		Port:            v.GetString("DB_OPERATION_PORT"),
		DatabaseName:    v.GetString("DB_OPERATION_DATABASE_NAME"),
		User:            v.GetString("DB_OPERATION_USER"),
		Password:        v.GetString("DB_OPERATION_PASSWORD"),
		SslMode:         v.GetString("DB_OPERATION_SSL_MODE"),
		Timezone:        v.GetString("DB_OPERATION_TIMEZONE"),
		MaxIdleConn:     v.GetInt("DB_OPERATION_MAX_IDLE_CONN"),
		MaxOpenConn:     v.GetInt("DB_OPERATION_MAX_OPEN_CONN"),
		MaxConnLifetime: time.Duration(v.GetInt("DB_OPERATION_MAX_CONN_LIFE_TIME")) * time.Second,
		MaxConnIdletime: time.Duration(v.GetInt("DB_OPERATION_MAX_CONN_IDLE_TIME")) * time.Second,
	}

	// // JWT Config
	// config.JWT = JwtConfig{
	// 	Expire:    v.GetInt64("JWT_EXPIRE"),
	// 	SecretKey: v.GetString("JWT_SECRET_KEY"),
	// }

	// Encryption Config
	config.Encrypt = EncryptConfig{
		KeySalt: v.GetString("KEY_SALT"),
	}

	// FTP Config
	config.FTP = FTPConfig{
		Host:     v.GetString("FTP_HOST"),
		Port:     v.GetInt64("FTP_PORT"),
		User:     v.GetString("FTP_USER"),
		Password: v.GetString("FTP_PASSWORD"),
		Path:     v.GetString("FTP_PATH"),
	}

	// File Local Path
	config.FileLocalPath = v.GetString("FILE_LOCAL_PATH")

	// Webhooks Config
	config.WebhooksConfig = WebhooksConfig{
		Customer: v.GetString("DISCORD_WEBHOOK_CUSTOMER"),
		Url:      v.GetString("DISCORD_WEBHOOK_URL"),
	}

	// Redis Configuration
	config.Redis = RedisConfig{
		Host:     v.GetString("REDIS_HOST"),
		Port:     v.GetInt("REDIS_PORT"),
		Password: v.GetString("REDIS_PASSWORD"),
		DB:       v.GetInt("REDIS_DB"),
	}

	// Set Service Mode
	config.ServiceMode = v.GetString("SERVICE_MODE")

	// Swagger Config
	config.Swagger = SwaggerConfig{
		SwaggerEnable: v.GetBool("SWAGGER_ENABLE"),
		SwaggerHost:   v.GetString("SWAGGER_HOST"),
		SwaggerScheme: v.GetString("SWAGGER_SCHEME"),
	}

	// Websocket Server and App Version
	config.WebSocketServer = v.GetString("WEBSOCKET_SERVER")
	config.VersionApp = v.GetString("VERSION_APP")

	// Google OAuth Config
	config.Google = GoogleConfig{
		ClientId: v.GetString("GOOGLE_CLIENT_ID"),
	}

	// Facebook OAuth Config
	config.Facebook = FacebookConfig{
		AppId:     v.GetString("FACEBOOK_APP_ID"),
		AppSecret: v.GetString("FACEBOOK_APP_SECRET"),
	}

	config.Apple = AppleConfig{
		ClientId: v.GetString("APPLE_CLIENT_ID"),
	}

	return nil
}

// Option define an option for config package
type Option func(*option)

// Get config
func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}
