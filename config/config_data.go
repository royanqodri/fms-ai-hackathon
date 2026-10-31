package config

import "time"

type Config struct {
	Service         ServiceConfig
	DatabaseConfig  DBConfig
	DBMainConfig    DBConfig
	JWT             JwtConfig
	Encrypt         EncryptConfig
	FTP             FTPConfig
	Redis           RedisConfig
	WebhooksConfig  WebhooksConfig
	WebSocket       WebSocketConfig
	WebSocketServer string
	API             APIConfig
	FileLocalPath   string
	VersionApp      string
	ServiceMode     string
	Swagger         SwaggerConfig
	Google          GoogleConfig
	Facebook        FacebookConfig
	Apple           AppleConfig
	Meta            MetaConfig
	TelegramToken   string
}

type ServiceConfig struct {
	Name         string
	Environment  string
	HttpAddress  string // include port, ex : 192.168.1.1:443
	LogLevel     string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	Location     string
}

type JwtConfig struct {
	Expire    int64
	SecretKey string
}

type EncryptConfig struct {
	KeySalt string
}

type DBConfig struct {
	Driver          string
	Host            string
	Port            string
	DatabaseName    string
	User            string
	Password        string
	SslMode         string
	Timezone        string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime time.Duration
	MaxConnIdletime time.Duration
}

type FTPConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	Path     string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type WebSocketConfig struct {
	Host string
	Port string
	Path string
}

type APIConfig struct {
	Host string
	Port string
}

type FirebaseConfig struct {
	FcmConfig FcmConfig
}

type FcmConfig struct {
	JsonPath  string
	ServerKey string
}

type WebhooksConfig struct {
	Customer string
	Url      string
}

type SwaggerConfig struct {
	SwaggerEnable bool
	SwaggerHost   string
	SwaggerScheme string
}

type GoogleConfig struct {
	ClientId string
}

type FacebookConfig struct {
	AppId     string
	AppSecret string
}

type AppleConfig struct {
	ClientId string // bundle ID (iOS) atau Services ID (web/android)
}

type MetaConfig struct {
	ClientId string
	ApiKey   string
}
