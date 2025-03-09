package admin

import (
	"context"
	"github.com/hashicorp/consul/api"
	"mydev/app/mydev/api/config"
	"mydev/app/pkg/options"
	gapp "mydev/gmicro/app"
	"mydev/gmicro/registry"
	"mydev/gmicro/registry/consul"
	"mydev/pkg/app"
	"mydev/pkg/log"
	"mydev/pkg/storage"
)

func NewApp(basename string) *app.App {
	cfg := config.New()
	appl := app.NewApp("api",
		"mydev",
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
		//不读配置 使用命令行参数时使用
		//app.WithNoConfig(),
	)
	return appl
}
func NewRegistrar(registry *options.RegistryOptions) registry.Registrar {
	c := api.DefaultConfig()
	c.Address = registry.Address
	c.Scheme = registry.Scheme
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}
func NewAPIApp(cfg *config.Config) (*gapp.App, error) {
	//初始化log
	log.Init(cfg.Log)
	defer log.Flush()
	//服务注册
	register := NewRegistrar(cfg.Registry)
	//连接redis
	redisConfig := &storage.Config{
		Host:                  cfg.Redis.Host,
		Port:                  cfg.Redis.Port,
		Addrs:                 cfg.Redis.Addrs,
		MasterName:            cfg.Redis.MasterName,
		Username:              cfg.Redis.Username,
		Password:              cfg.Redis.Password,
		Database:              cfg.Redis.Database,
		MaxIdle:               cfg.Redis.MaxIdle,
		MaxActive:             cfg.Redis.MaxActive,
		Timeout:               cfg.Redis.Timeout,
		EnableCluster:         cfg.Redis.EnableCluster,
		UseSSL:                cfg.Redis.UseSSL,
		SSLInsecureSkipVerify: cfg.Redis.SSLInsecureSkipVerify,
		EnableTracing:         cfg.Redis.EnableTracing,
	}
	go storage.ConnectToRedis(context.Background(), redisConfig)
	//生成http服务
	restServer, err := NewAPIHTTPServer(cfg)
	if err != nil {
		return nil, err
	}
	return gapp.New(
		gapp.WithName(cfg.Server.Name),
		gapp.WithRestServer(restServer),
		gapp.WithRegistrar(register),
	), nil
}

// controller(参数校验) ->service(具体的业务逻辑)->(数据库的接口)
func run(cfg *config.Config) app.RunFunc {
	return func(baseName string) error {
		apiApp, err := NewAPIApp(cfg)
		if err != nil {
			return err
		}
		//启动
		if err := apiApp.Run(); err != nil {
			log.Errorf("run api app error: %s", err)
		}
		return nil
	}
}

/*
现在我想换一个rpc,zrpc
逻辑和rpc的数据耦合了数据层,我们和gorm耦合了我们很有可能会遇到一下方面的问题:
rpc我们可能会换
底层orm可能会面临两个问题: 1。我们想优化性能，我想优化goods列表的查询性能,我们查询现在是从es中查询的,后面我们想要从hbase中进行查询
web：
用了gin,我们想换了krqtosihttp go-zero的http服务
注册中心想换，consul我们想换成nacos，我们想换成k8s的服务发现和注册
缓存想换redis，后面我们可能想使用内存，memcache
*/
