package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	v1 "mydev/api/goods/v1"
	"mydev/gmicro/registry/consul"
	rpc "mydev/gmicro/server/rpcserver"
	_ "mydev/gmicro/server/rpcserver/resolver/direct" // 这个是直接连接的 下面已经实现watcher长轮询了  弃用
	"mydev/gmicro/server/rpcserver/selector"
	"mydev/gmicro/server/rpcserver/selector/random"
)

func main() {
	//设置全局的负载均衡策略
	selector.SetGlobalSelector(random.NewBuilder())
	rpc.InitBuilder()

	conf := api.DefaultConfig()
	conf.Address = "127.0.0.1:8500"
	conf.Scheme = "http"
	cli, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	conn, err := rpc.DialInsecure(context.Background(),
		rpc.WithBalancerName("selector"),
		rpc.WithDiscovery(r),
		/*
			第3个/是为了第二个参数是空的
			默认格式：direct://<authority>/127.0.0.1:8078
			以后使用nacos或者其他的中心 也不用改discovery 只修改conf就可以
			服务发现可以直接去kartors里面copy registry下的etcd nacos等使用
		*/
		rpc.WithEndpoint("discovery:///mydev-goods-srv"),
		//rpc.WithClientTimeout(time.Second*1),
	)
	//conn, err := grpc.Dial("127.0.0.1:8078", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	uc := v1.NewGoodsClient(conn)
	re, err := uc.GoodsList(context.Background(), &v1.GoodsFilterRequest{
		KeyWords: "西红柿",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(re)
}
