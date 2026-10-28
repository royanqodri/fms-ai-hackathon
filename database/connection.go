package database

import (
	"fmt"
	"time"

	"github.com/royanqodri/fms-ai-hackathon/config"
	"github.com/royanqodri/fms-ai-hackathon/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

func ConnectDatabaseMaster() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Get().DBMainConfig.Host,
		config.Get().DBMainConfig.User,
		config.Get().DBMainConfig.Password,
		config.Get().DBMainConfig.DatabaseName,
		config.Get().DBMainConfig.Port,
		config.Get().DBMainConfig.SslMode,
		config.Get().DBMainConfig.Timezone,
	)
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	util.PanicIfError(err)

	sqlDB, err := DBConn.DB()
	util.PanicIfError(err)

	sqlDB.SetMaxIdleConns(config.Get().DBMainConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.Get().DBMainConfig.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(config.Get().DBMainConfig.MaxConnLifetime * time.Minute)
	sqlDB.SetConnMaxIdleTime(config.Get().DBMainConfig.MaxConnIdletime * time.Minute)

	sqlDB.Exec("SET statement_timeout = '15s'")
}
