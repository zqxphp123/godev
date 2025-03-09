package service

import (
	"mydev/app/mydev/api/internal/data"
	goods2 "mydev/app/mydev/api/internal/service/goods/v1"
	v1 "mydev/app/mydev/api/internal/service/sms/v1"
	user2 "mydev/app/mydev/api/internal/service/user/v1"
	"mydev/app/pkg/options"
)

type ServiceFactory interface {
	Goods() goods2.GoodsSrv
	Users() user2.UserSrv
	Base() v1.SmsSrv
}

// 持有化
type service struct {
	data    data.DataFactory
	jwtOpts *options.JwtOptions
	smsOpts *options.SmsOptions
}

// 这里不用sevice是为了防止循环引用，其实也可以把service放到同一个文件夹下，但是那样文件多了就分不清了  不好看

func (s *service) Users() user2.UserSrv {
	return user2.NewUser(s.data, s.jwtOpts)
}

func (s *service) Base() v1.SmsSrv {
	return v1.NewSms(s.smsOpts)
}

func (s *service) Goods() goods2.GoodsSrv {
	return goods2.NewGoods(s.data)
}

func NewService(data data.DataFactory, jwtOpts *options.JwtOptions, smsOpts *options.SmsOptions) *service {
	return &service{
		data:    data,
		jwtOpts: jwtOpts,
		smsOpts: smsOpts,
	}
}

var _ ServiceFactory = &service{}
