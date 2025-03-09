package rpc

import (
	"fmt"
	consulAPI "github.com/hashicorp/consul/api"
	"mydev/app/pkg/code"
	"mydev/pkg/errors"
	"sync"

	gpb "mydev/api/goods/v1"
	upb "mydev/api/user/v1"
	"mydev/app/mydev/api/internal/data"
	"mydev/app/pkg/options"
	"mydev/gmicro/registry"
	"mydev/gmicro/registry/consul"
)

type grpcData struct {
	gc gpb.GoodsClient
	uc upb.UserClient
}

func (g *grpcData) Goods() gpb.GoodsClient {
	return g.gc
}

func (g *grpcData) Users() data.UserData {
	return NewUsers(g.uc)
}

var (
	rpcFactory data.DataFactory
	once       sync.Once
)

// 目前是基于consul实现的  以后想换成nocos etcd等  可以直接在这换
func NewDiscovery(opts *options.RegistryOptions) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = opts.Address
	c.Scheme = opts.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}

// rpc的连接，基于服务发现
func GetDataFactoryOr(options *options.RegistryOptions) (data.DataFactory, error) {
	if options == nil && rpcFactory == nil {
		return nil, fmt.Errorf("failed to get rpc store fatory")
	}
	once.Do(func() {
		discovery := NewDiscovery(options)
		//这里负责依赖所有的rpc连接
		userClient := NewUserServiceClient(discovery)
		goodsClient := NewGoodsServiceClient(discovery)
		rpcFactory = &grpcData{
			gc: goodsClient,
			uc: userClient,
		}
	})
	var err error
	if rpcFactory == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectGRPC, "failed to get rpc store factory")
	}

	return rpcFactory, nil
}

//func NewUserClient()
