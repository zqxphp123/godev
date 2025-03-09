package v1

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	v1 "mydev/app/inventory/srv/internal/data/v1"
	"mydev/app/pkg/options"
)

type ServiceFactory interface {
	Inventorys() InventorySrv
}
type service struct {
	data      v1.DataFactory
	redisOpts *options.RedisOptions
	pool      redsyncredis.Pool
}

func (s *service) Inventorys() InventorySrv {
	return newInventoryService(s)
}

func NewService(data v1.DataFactory, redisOpts *options.RedisOptions) *service {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", redisOpts.Host, redisOpts.Port),
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	return &service{data: data, redisOpts: redisOpts, pool: pool}
}

var _ ServiceFactory = &service{}
