package user

import (
	"context"

	upbv1 "mydev/api/user/v1"

	"mydev/pkg/log"
)

func (u *userServer) GetUserByMobile(ctx context.Context, request *upbv1.MobileRequest) (*upbv1.UserInfoResponse, error) {
	log.Info("get user by mobile function called.")
	user, err := u.srv.GetByMobile(ctx, request.Mobile)
	if err != nil {
		log.Errorf("get user by mobile: %s,error: %v", request.Mobile, err)
	}
	userInfoRsp := DTOToResponse(*user)
	return userInfoRsp, nil
}
