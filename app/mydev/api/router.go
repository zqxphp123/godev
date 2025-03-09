package admin

import (
	"mydev/app/mydev/api/config"
	"mydev/app/mydev/api/internal/controller/goods/v1"
	sms2 "mydev/app/mydev/api/internal/controller/sms/v1"
	"mydev/app/mydev/api/internal/controller/user/v1"
	"mydev/app/mydev/api/internal/data/rpc"
	"mydev/app/mydev/api/internal/service"
	"mydev/gmicro/server/restserver"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")
	ugroup := v1.Group("/user")
	baseRouter := v1.Group("/base")
	goodsRouter := v1.Group("/goods")

	// rpc的连接，基于服务发现
	data, err := rpc.GetDataFactoryOr(cfg.Registry)
	if err != nil {
		panic(err)
	}
	serviceFactory := service.NewService(data, cfg.Jwt, cfg.Sms)

	//用户相关的api
	uController := user.NewUserController(g.Translator(), serviceFactory)
	{
		ugroup.POST("pwd_login", uController.Login)
		ugroup.POST("register", uController.Regisger)

		jwtAuth := newJWTAuth(cfg.Jwt)
		//第三方登录模式 暂不用
		//jwtStragy:=jwtAuth.(auth.JWTStrategy)
		//jwtStragy.LoginHandler()
		ugroup.GET("detail", jwtAuth.AuthFunc(), uController.GetUserDetail)
		ugroup.PATCH("update", jwtAuth.AuthFunc(), uController.UpdateUser)
	}
	//验证相关的api
	smsCtl := sms2.NewSmsController(serviceFactory, g.Translator())
	{
		baseRouter.POST("send_sms", smsCtl.SendSms)
		baseRouter.GET("captcha", user.GetCaptcha)
	}
	//商品相关的api
	{
		goodsController := goods.NewGoodsController(serviceFactory, g.Translator())
		goodsRouter.GET("", goodsController.List)
	}
}
