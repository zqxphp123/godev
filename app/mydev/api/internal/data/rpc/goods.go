package rpc

import (
	"context"

	gpbv1 "mydev/api/goods/v1"
	"mydev/gmicro/registry"
	"mydev/gmicro/server/rpcserver"
	"mydev/gmicro/server/rpcserver/clientinterceptors"
)

const goodsServiceName = "discovery:///mydev-goods-srv"

func NewGoodsServiceClient(r registry.Discovery) gpbv1.GoodsClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(goodsServiceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := gpbv1.NewGoodsClient(conn)
	return c
}
