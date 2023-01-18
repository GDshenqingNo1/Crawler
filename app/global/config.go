package global

import (
	"WebCrawler/app/internal/model/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	MysqlDB *gorm.DB
	Config  *config.Config
	Logger  *zap.Logger
)
