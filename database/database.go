package database

import (
	"github.com/goplaceapp/goplace-user/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Database struct {
	GetSharedDB func() *gorm.DB
	GetTenantDB func(tenant string) *gorm.DB
}

var gormLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  utils.GetLogLevel(),
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
		Colorful:                  true,
	},
)
