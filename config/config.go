package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

var (
	config *Config
)

type option struct {
}

func Init() error {
	v := viper.New()
	v.AutomaticEnv()

	mode := v.GetString("SERVICE_ENV")
	if mode == "" {
		mode = "LOCAL"
	}

	if mode != "PROD" {
		v.SetConfigName("config")
		v.SetConfigType("env")
		v.AddConfigPath(".")
		v.AddConfigPath("..") // Tambahkan ini agar bisa dibaca dari folder cmd
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	v.SetDefault("SERVICE_NAME", "fms-ai-hackathon")
	v.SetDefault("VERSION_APP", "latest")

	config = &Config{}

	// --- SERVICE CONFIG ---
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

	// --- DB CONFIG (1 DB SAJA) ---
	// Kita isi DBMasterConfig dengan variabel DB_...
	config.DBMainConfig = DBConfig{
		Driver:          v.GetString("DB_DRIVER"),
		Host:            v.GetString("DB_HOST"),
		Port:            v.GetString("DB_PORT"),
		DatabaseName:    v.GetString("DB_NAME"),
		User:            v.GetString("DB_USER"),
		Password:        v.GetString("DB_PASSWORD"),
		SslMode:         v.GetString("DB_SSL_MODE"),
		Timezone:        v.GetString("DB_TIMEZONE"),
		MaxIdleConn:     v.GetInt("DB_MAX_IDLE_CONN"),
		MaxOpenConn:     v.GetInt("DB_MAX_OPEN_CONN"),
		MaxConnLifetime: time.Duration(v.GetInt("DB_MAX_CONN_LIFE_TIME")) * time.Second,
		MaxConnIdletime: time.Duration(v.GetInt("DB_MAX_CONN_IDLE_TIME")) * time.Second,
	}

	// --- KONFIGURASI LAIN (Dikosongkan agar struct tetap aman) ---
	config.JWT = JwtConfig{
		Expire:    v.GetInt64("JWT_EXPIRE"),
		SecretKey: v.GetString("JWT_SECRET_KEY"),
	}
	config.Encrypt = EncryptConfig{KeySalt: v.GetString("KEY_SALT")}
	config.FTP = FTPConfig{
		Host:     v.GetString("FTP_HOST"),
		Port:     v.GetInt64("FTP_PORT"),
		User:     v.GetString("FTP_USER"),
		Password: v.GetString("FTP_PASSWORD"),
		Path:     v.GetString("FTP_PATH"),
	}
	config.FileLocalPath = v.GetString("FILE_LOCAL_PATH")
	config.WebhooksConfig = WebhooksConfig{
		Customer: v.GetString("DISCORD_WEBHOOK_CUSTOMER"),
		Url:      v.GetString("DISCORD_WEBHOOK_URL"),
	}
	config.Redis = RedisConfig{
		Host:     v.GetString("REDIS_HOST"),
		Port:     v.GetInt("REDIS_PORT"),
		Password: v.GetString("REDIS_PASSWORD"),
		DB:       v.GetInt("REDIS_DB"),
	}
	config.ServiceMode = v.GetString("SERVICE_MODE")
	config.Swagger = SwaggerConfig{
		SwaggerEnable: v.GetBool("SWAGGER_ENABLE"),
		SwaggerHost:   v.GetString("SWAGGER_HOST"),
		SwaggerScheme: v.GetString("SWAGGER_SCHEME"),
	}
	config.WebSocketServer = v.GetString("WEBSOCKET_SERVER")
	config.VersionApp = v.GetString("VERSION_APP")
	config.Google = GoogleConfig{ClientId: v.GetString("GOOGLE_CLIENT_ID")}
	config.Facebook = FacebookConfig{
		AppId:     v.GetString("FACEBOOK_APP_ID"),
		AppSecret: v.GetString("FACEBOOK_APP_SECRET"),
	}
	config.Apple = AppleConfig{ClientId: v.GetString("APPLE_CLIENT_ID")}

	config.Meta = MetaConfig{
		ApiKey: os.Getenv("META_API_KEY"),
	}

	config.TelegramToken = os.Getenv("TELEGRAM_BOT_TOKEN")

	return nil
}

func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}
