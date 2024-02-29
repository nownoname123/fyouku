package dao

import (
	"fyoukuapi/config"
	"fyoukuapi/pkg/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var (
	Db  *gorm.DB
	err error
)

func init() {
	Db, err = gorm.Open("mysql", config.MysqlDb)
	Db.SingularTable(true)
	if err != nil {
		logger.Error(map[string]interface{}{"mysql connect error": err.Error()})
	}
	if Db.Error != nil {
		logger.Error(map[string]interface{}{"database error": Db.Error})
	}

	// ----------------------- 连接池设置 -----------------------
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	Db.DB().SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	Db.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	Db.DB().SetConnMaxLifetime(time.Hour)
}
