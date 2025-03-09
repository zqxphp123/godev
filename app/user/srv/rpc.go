package srv

import (
	"fmt"
	upbv1 "mydev/api/user/v1"
	"mydev/app/user/srv/config"
	"mydev/app/user/srv/controller/user"
	"mydev/app/user/srv/data/v1/db"
	srv1 "mydev/app/user/srv/service/v1"
	"mydev/gmicro/core/trace"
	"mydev/gmicro/server/rpcserver"
	"mydev/pkg/log"
)

func NewUserRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	//初始化open-telemetry的exporter
	trace.InitAgent(trace.Options{
		cfg.Telemetry.Name,
		cfg.Telemetry.Endpoint,
		cfg.Telemetry.Sampler,
		cfg.Telemetry.Batcher,
	})
	//data := mock.NewUsers()
	gormDB, err := db.GetDBFactoryOr(cfg.MySQLOptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	data := db.NewUsers(gormDB)
	srv := srv1.NewuserService(data)
	userver := user.NewUserServer(srv)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	urpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))
	upbv1.RegisterUserServer(urpcServer.Server, userver)

	return urpcServer, nil

}
