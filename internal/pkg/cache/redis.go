package cache

import (
	"at-kratos/internal/conf"
	"github.com/go-redis/redis"
	"log"
)

var rdb *redis.Client

// Init Redis
func Init(conf *conf.Data_Redis) *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       int(conf.DbNo),
		PoolSize: int(conf.PoolSize),
	})
	if _, err := rdb.Ping().Result(); err != nil {
		log.Fatal(err)
	}
	return rdb
}
