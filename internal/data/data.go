package data

import (
	"at-kratos/internal/conf"
	"at-kratos/internal/pkg/cache"
	"at-kratos/internal/pkg/database"
	"gitee.com/chunanyong/zorm"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewSignupLogRepo,
)

// Data .
type Data struct {
	db    *zorm.DBDao
	redis *redis.Client
}

// NewData .
func NewData(conf *conf.Data, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(logger)
	d := new(Data)
	d.db = database.Init(conf.Database)
	d.redis = cache.Init(conf.Redis)
	return d, func() {
		logHelper.Info("message", "closing the data resources")
	}, nil
}
