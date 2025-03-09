package admin

import (
	"mydev/app/mydev/api/config"
	"mydev/gmicro/server/restserver"
)

func NewAPIHTTPServer(cfg *config.Config) (*restserver.Server, error) {

	aRestServer := restserver.NewServer(
		restserver.WithPort(cfg.Server.HttpPort),
		restserver.WithMiddlewares(cfg.Server.Middlewares),
		restserver.WithEnableProfiling(true),
		restserver.WithMetrics(true),
	)
	//配置好路由
	initRouter(aRestServer, cfg)
	return aRestServer, nil
}
