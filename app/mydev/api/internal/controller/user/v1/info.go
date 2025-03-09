package user

import (
	"github.com/gin-gonic/gin"
	"mydev/gmicro/server/restserver/middlewares"
	"mydev/pkg/common/core"
	"time"
)

func (us *userServer) GetUserDetail(ctx *gin.Context) {
	userID, _ := ctx.Get(middlewares.KeyUserID)
	//最好debug看一下是什么类型  再进行断言
	userDTO, err := us.srv.Users().Get(ctx, uint64(userID.(float64)))
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	core.WriteResponse(ctx, nil, gin.H{
		"name":     userDTO.NickName,
		"birthday": userDTO.Birthday.Format(time.DateOnly),
		"gender":   userDTO.Gender,
		"mobile":   userDTO.Mobile,
	})

}
