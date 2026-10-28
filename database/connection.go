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
	DBConnMaster    *gorm.DB
	DBConnOperation *gorm.DB
)

func ConnectDatabaseMaster() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Get().DBMasterConfig.Host,
		config.Get().DBMasterConfig.User,
		config.Get().DBMasterConfig.Password,
		config.Get().DBMasterConfig.DatabaseName,
		config.Get().DBMasterConfig.Port,
		config.Get().DBMasterConfig.SslMode,
		config.Get().DBMasterConfig.Timezone,
	)
	DBConnMaster, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	util.PanicIfError(err)

	sqlDB, err := DBConnMaster.DB()
	util.PanicIfError(err)

	sqlDB.SetMaxIdleConns(config.Get().DBMasterConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.Get().DBMasterConfig.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(config.Get().DBMasterConfig.MaxConnLifetime * time.Minute)
	sqlDB.SetConnMaxIdleTime(config.Get().DBMasterConfig.MaxConnIdletime * time.Minute)

	sqlDB.Exec("SET statement_timeout = '15s'")
}

func ConnectDatabaseOperation() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Get().DBOperationConfig.Host,
		config.Get().DBOperationConfig.User,
		config.Get().DBOperationConfig.Password,
		config.Get().DBOperationConfig.DatabaseName,
		config.Get().DBOperationConfig.Port,
		config.Get().DBOperationConfig.SslMode,
		config.Get().DBOperationConfig.Timezone,
	)
	DBConnOperation, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	util.PanicIfError(err)

	sqlDB, err := DBConnOperation.DB()
	util.PanicIfError(err)

	if err = sqlDB.Ping(); err != nil {
		util.PanicIfError(err)
	}
	sqlDB.SetMaxIdleConns(config.Get().DBOperationConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.Get().DBOperationConfig.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(config.Get().DBOperationConfig.MaxConnLifetime * time.Minute)
	sqlDB.SetConnMaxIdleTime(config.Get().DBOperationConfig.MaxConnIdletime * time.Minute)

	sqlDB.Exec("SET statement_timeout = '15s'")
}
