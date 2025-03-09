package srv

import (
	"fmt"
	gpb "mydev/api/order/v1"
	"mydev/app/order/srv/config"
	"mydev/app/order/srv/internal/controller/order/v1"
	db2 "mydev/app/order/srv/internal/data/v1/db"
	v13 "mydev/app/order/srv/internal/service/v1"
	"mydev/gmicro/core/trace"
	"mydev/gmicro/server/rpcserver"

	"mydev/pkg/log"
)

func NewOrderRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	//初始化open-telemetry的exporter
	trace.InitAgent(trace.Options{
		cfg.Telemetry.Name,
		cfg.Telemetry.Endpoint,
		cfg.Telemetry.Sampler,
		cfg.Telemetry.Batcher,
	})

	dataFactory, err := db2.GetDataFactoryOr(cfg.MySQLOptions, cfg.Registry)
	if err != nil {
		log.Fatal(err.Error())
	}

	orderSrvFactory := v13.NewService(dataFactory, cfg.Dtm)
	orderServer := order.NewOrderServer(orderSrvFactory)
	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))
	gpb.RegisterOrderServer(grpcServer.Server, orderServer)
	return grpcServer, nil
}
