package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"math/rand"
	v1 "mydev/api/order/v1"
	"mydev/gmicro/registry/consul"
	rpc "mydev/gmicro/server/rpcserver"
	_ "mydev/gmicro/server/rpcserver/resolver/direct" // 这个是直接连接的 下面已经实现watcher长轮询了  弃用
	"mydev/gmicro/server/rpcserver/selector"
	"mydev/gmicro/server/rpcserver/selector/random"
	"time"
)

func generateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}
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
		rpc.WithEndpoint("discovery:///mydev-order-srv"),
		//rpc.WithClientTimeout(time.Second*1),
	)
	//conn, err := grpc.Dial("127.0.0.1:8078", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	uc := v1.NewOrderClient(conn)
	_, err = uc.SubmitOrder(context.Background(), &v1.OrderRequest{
		UserId:  12,
		Address: "山东商业职业技术学院",
		OrderSn: generateOrderSn(12),
		Name:    "jzin",
		Post:    "尽快发货",
		Mobile:  "15325098743",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("订单新建成功")
}
