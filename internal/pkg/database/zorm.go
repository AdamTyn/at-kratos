package database

import (
	"at-kratos/internal/conf"
	"gitee.com/chunanyong/zorm"
	_ "github.com/lib/pq"
)

// Init Postgresql
func Init(conf *conf.Data_Database) *zorm.DBDao {
	dbConfig := zorm.DataSourceConfig{
		DSN:                   conf.Source,
		DriverName:            conf.Driver,
		DBType:                conf.Type,
		ConnMaxLifetimeSecond: 600,
		PrintSQL:              true,
	}
	db, err := zorm.NewDBDao(&dbConfig)
	if err != nil {
		panic(err)
	}
	return db
}
