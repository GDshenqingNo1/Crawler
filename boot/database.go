package boot

import (
	g "WebCrawler/app/global"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlDBSetup() {
	config := g.Config.Database.Mysql
	db, err := gorm.Open(mysql.Open(config.GetDsn()))
	if err != nil {
		g.Logger.Fatal("initialize mysql failed", zap.Error(err))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(g.Config.Database.Mysql.GetConnMaxIdleTime())
	sqlDB.SetConnMaxLifetime(g.Config.Database.Mysql.GetConnMaxIdleTime())
	sqlDB.SetMaxIdleConns(g.Config.Database.Mysql.MaxIdleConns)
	sqlDB.SetMaxIdleConns(g.Config.Database.Mysql.MaxOpenConns)
	err = sqlDB.Ping()
	if err != nil {
		g.Logger.Fatal("connect to mysql db failed.", zap.Error(err))
	}
	g.MysqlDB = db
	g.Logger.Info("initialize mysql successfully.")
}
