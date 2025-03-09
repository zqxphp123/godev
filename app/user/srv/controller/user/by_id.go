package user

import (
	"context"

	upbv1 "mydev/api/user/v1"

	"mydev/pkg/log"
)

func (u *userServer) GetUserById(ctx context.Context, request *upbv1.IdRequest) (*upbv1.UserInfoResponse, error) {
	log.Info("get user by id function called.")
	user, err := u.srv.GetByID(ctx, uint64(request.Id))
	if err != nil {
		log.Errorf("get user by id: %s,error: %v", request.Id, err)
	}
	userInfoRsp := DTOToResponse(*user)
	return userInfoRsp, nil
}
