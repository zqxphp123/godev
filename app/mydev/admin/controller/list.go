package controller

import (
	"github.com/gin-gonic/gin"
	"mydev/pkg/log"
)

func (us *userServer) List(ctx *gin.Context) {
	log.Info("GetUserList is called")
}
