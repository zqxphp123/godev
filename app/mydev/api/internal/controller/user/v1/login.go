package user

import (
	"github.com/gin-gonic/gin"
	gin2 "mydev/app/pkg/translator/gin"
	"mydev/pkg/log"
	"net/http"
)

type PassWordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=10"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

func (us *userServer) Login(ctx *gin.Context) {
	log.Info("login is called")
	//表单验证
	passwordLoginForm := PassWordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		gin2.HandleValidatorError(ctx, err, us.trans)
		return
	}
	//验证码验证
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		//return
	}
	//登录验证
	userDTO, err := us.srv.Users().MobileLogin(ctx, passwordLoginForm.Mobile, passwordLoginForm.PassWord)
	if err != nil {
		//core.WriteResponse(ctx, http.StatusInternalServerError, "登陆失败", nil)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":          userDTO.ID,
		"nick_name":   userDTO.NickName,
		"token":       userDTO.Token,
		"expiresd_at": userDTO.ExpiresAt,
	})

}
