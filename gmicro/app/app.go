package app

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"mydev/gmicro/registry"
	gs "mydev/gmicro/server"
	"mydev/pkg/log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	opts options
	lk   sync.Mutex
	//受保护
	instance *registry.ServiceInstance

	cancel func()
}

func New(opts ...Option) *App {
	o := options{
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
		stopTimeout:      10 * time.Second,
	}

	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &App{
		opts: o,
	}
}

// 启动整个服务
func (a *App) Run() error {
	//注册的信息
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	//这个变量可能被其他的goroutine访问
	a.lk.Lock()
	a.instance = instance
	a.lk.Unlock()

	//重点 写的很简单，http服务要启动

	//go func() {
	//	if a.opts.rpcServer != nil {
	//		err := a.opts.rpcServer.Start()
	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	//}()

	//现在启动了两个server，一个是http，一个是rpc
	/*
		这两个server是否必须启动成功？
		如果有一个启动失败，那么我们就要停止另外一个server
		启动了多个，如果其中一个启动失败，其他应该被取消
		如果剩余的server的状态：
			1.没有开始调用start
				不进行就行或者调用stop
			2.start进行中
				调用进行中的cancel
			3.start已经完成
				调用stop
	*/

	var servers []gs.Server
	if a.opts.restServer != nil {
		servers = append(servers, a.opts.restServer)
	}
	if a.opts.rpcServer != nil {
		servers = append(servers, a.opts.rpcServer)
	}
	//app在stop的时候想要通知到服务下进行cancel
	//这时候我们自己生成一个context，把cancel方法注入到app当中，这时候在stop的时候cancel方法就能通知到服务中

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	for _, srv := range servers {
		//启动server
		//再启动一个goroutine 去监听是否有err产生
		//goroutine在for循环中的坑↓
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() //wait for stop signal
			//不可能无休止的等待stop信号
			sctx, cancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(sctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			log.Info("start rest server")
			//context的作用：应该可以接受一个可以cancel的context 随时取消
			return srv.Start(ctx)
		})
	}
	wg.Wait()

	//err = eg.Wait()
	//if err != nil {
	//}

	//注册服务
	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		defer rcancel()
		err = a.opts.registrar.Register(rctx, instance)
		if err != nil {
			log.Errorf("registrar service error: %s", err)
			return err
		}
	}
	//监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	//<-c
	//由于a.cancel()执行的很快 导致整个goroutine程序退出  所以放到goroutine里监听chan。达到一个阻塞的效果
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return a.Stop()
		}
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// 停止服务
func (a *App) Stop() error {
	a.lk.Lock()
	instance := a.instance
	a.lk.Unlock()

	log.Info("start deregister service")
	if a.opts.registrar != nil && instance != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer rcancel()
		err := a.opts.registrar.Deregister(rctx, instance)
		if err != nil {
			log.Errorf("deregister service error: %s", err)
			return err
		}
	}
	//自己生成的context生成cancel后往服务中传递，所以能通知到所有的服务下的context
	if a.cancel != nil {
		log.Infof("start cancel context")
		a.cancel()
	}
	return nil
}

// 创建服务注册结构体
func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	//从rpcserver，restserver去主动获取这些信息
	if a.opts.rpcServer != nil {
		if a.opts.rpcServer.Endpoint() != nil {
			endpoints = append(endpoints, a.opts.rpcServer.Endpoint().String())
		} else {
			u := &url.URL{
				Scheme: "grpc",
				Host:   a.opts.rpcServer.Address(),
			}
			endpoints = append(endpoints, u.String())
		}
	}
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endpoints,
	}, nil
}
