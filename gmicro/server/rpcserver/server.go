package rpcserver

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	apimd "mydev/api/metadata"
	srvintc "mydev/gmicro/server/rpcserver/serverinterceptors"
	"mydev/pkg/host"
	"mydev/pkg/log"
	"net"
	"net/url"
	"time"
)

type ServerOption func(o *Server)
type Server struct {
	*grpc.Server

	address string
	//可传递拦截器
	unaryInts []grpc.UnaryServerInterceptor
	//stream
	streamInts []grpc.StreamServerInterceptor
	//客户自己实现ServerOption
	grpcOpts []grpc.ServerOption
	//监听器
	lis net.Listener
	//timeout
	timeout time.Duration

	//rpc的健康检查
	health *health.Server
	//一个grpc接口查看所有rpc服务
	metadata *apimd.Server
	//设置监听的ip port
	endpoint *url.URL
	//由用户决定是否加入指标
	enableMetrics bool
}

func (s *Server) Endpoint() *url.URL {
	return s.endpoint
}
func (s *Server) Address() string {
	return s.address
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		//自己生成端口号
		address: ":0",
		health:  health.NewServer(),
		//timeout: 1 * time.Second,
	}
	for _, o := range opts {
		o(srv)
	}
	//我们现在希望用户不设置拦截器的情况下，我们会自动默认加上一些必须的拦截器 , crash tracing
	unaryInts := []grpc.UnaryServerInterceptor{
		srvintc.UnaryCrashInterceptor,
		//闭包特性，直接调用后返回resp供grpc拦截器调用
		otelgrpc.UnaryServerInterceptor(),
		//srvintc.UnaryTimeoutInterceptor(srv.timeout),
	}
	//由用户决定是否加入指标
	if srv.enableMetrics {
		unaryInts = append(unaryInts, srvintc.UnaryPrometheusInterceptor)
	}
	//timeout可以交给用户设置，不设置就不用此拦截器
	if srv.timeout > 0 {
		unaryInts = append(unaryInts, srvintc.UnaryTimeoutInterceptor(srv.timeout))
	}
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}

	//把我们传入的拦截器转换成grpc的ServerOption
	grpcOpts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(unaryInts...)}
	//把用户自己传入的grpc.ServerOption放在一起
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)
	//注册metadata的Server
	srv.metadata = apimd.NewServer(srv.Server)
	//自动解析address
	err := srv.listenAndEndpoint()
	if err != nil {
		return nil
	}
	//注册health
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	//可以支持用户直接通过grpc的一个接口查看当前支持的所有的rpc服务
	apimd.RegisterMetadataServer(srv.Server, srv.metadata)
	reflection.Register(srv.Server)
	return srv
}

func WithAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}
func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}
func WithLis(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}
func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInts = in
	}
}
func WithStreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInts = in
	}
}
func WithOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}
func WithenableMetrics(metrics bool) ServerOption {
	return func(s *Server) {
		s.enableMetrics = metrics
	}
}

// 完成ip和端口的提取
func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	return nil
}

// 完成ip和port的提取
//
//	func (s *Server) listenAndEndpoint() error {
//		if s.lis == nil {
//			lis, err := net.Listen("tpc", s.address)
//			if err != nil {
//				return err
//			}
//			s.lis = lis
//		}
//		addr, err := host.Extract(s.address, s.lis)
//		if err != nil {
//			_ = s.lis.Close()
//			return err
//		}
//		s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
//		return nil
//	}
func (s *Server) Start(ctx context.Context) error {
	log.Infof("[grpc] server listening on: %s", s.lis.Addr().String())
	//改grpc核心变量 状态
	//只有.Resume()之后，请求才能进来
	//s.health.Shutdown()相反
	s.health.Resume()
	return s.Server.Serve(s.lis)

}
func (s *Server) Stop(ctx context.Context) error {
	//设置服务的状态为not_serving 防止接受新的请求
	s.health.Shutdown()
	s.GracefulStop()
	log.Infof("[grpc] server stopped")
	return nil
}
