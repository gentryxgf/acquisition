package redis

import (
	"context"
	"fmt"
	"server/config"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	Rs *redsync.Redsync
	Rc *redis.Client
)

func Init(cfg *config.RedisConfig) error {
	Rc = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})

	err := Rc.Ping(context.Background()).Err()
	if err != nil {
		return err
	}

	zap.L().Info("Init Redis Success!")

	pool := goredis.NewPool(Rc)

	Rs = redsync.New(pool)

	zap.L().Info("Create redisync instance success!")
	return nil
}
