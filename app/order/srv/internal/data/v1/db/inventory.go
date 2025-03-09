package db

import (
	"context"
	proto "mydev/api/inventory/v1"
	"mydev/app/pkg/options"
	"mydev/gmicro/registry"
	"mydev/gmicro/server/rpcserver"
	"mydev/gmicro/server/rpcserver/clientinterceptors"
)

const invServiceName = "discovery:///mydev-inventory-srv"

func GetInventoryClient(opts *options.RegistryOptions) proto.InventoryClient {
	discovery := NewDiscovery(opts)
	//这里负责依赖所有的rpc连接
	inventoryClient := NewInventoryServiceClient(discovery)
	return inventoryClient
}

func NewInventoryServiceClient(r registry.Discovery) proto.InventoryClient {

	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(invServiceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := proto.NewInventoryClient(conn)
	return c
}
