package app

import (
	"mydev/gmicro/registry"
	"mydev/gmicro/server/restserver"
	"mydev/gmicro/server/rpcserver"
	"net/url"
	"os"
	"time"
)

type Option func(o *options)
type options struct {
	id        string
	name      string
	endpoints []*url.URL
	sigs      []os.Signal
	//允许用户传入自己的实现
	registrarTimeout time.Duration
	registrar        registry.Registrar

	//stop超时时间
	stopTimeout time.Duration

	//传递http服务
	restServer *restserver.Server
	//传递rpc服务
	rpcServer *rpcserver.Server
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithRPCServer(server *rpcserver.Server) Option {
	return func(o *options) {
		o.rpcServer = server
	}
}
func WithRestServer(server *restserver.Server) Option {
	return func(o *options) {
		o.restServer = server
	}
}

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}
func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}
func WithSigs(sigs []os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}
func WithRegistrarTimeout(registrarTimeout time.Duration) Option {
	return func(o *options) {
		o.registrarTimeout = registrarTimeout
	}
}
