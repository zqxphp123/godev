package rpc

import (
	"context"
	upbv1 "mydev/api/user/v1"
	"mydev/app/mydev/api/internal/data"
	"mydev/app/pkg/code"
	"mydev/gmicro/registry"
	"mydev/gmicro/server/rpcserver"
	"mydev/gmicro/server/rpcserver/clientinterceptors"
	itime "mydev/pkg/common/time"
	"mydev/pkg/errors"
	"time"
)

const userServiceName = "discovery:///mydev-user-srv"

type users struct {
	uc upbv1.UserClient
}

func NewUserServiceClient(r registry.Discovery) upbv1.UserClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(userServiceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := upbv1.NewUserClient(conn)
	return c
}
func NewUsers(uc upbv1.UserClient) *users {
	return &users{uc}
}
func (u *users) Create(ctx context.Context, user *data.User) error {
	protoUser := &upbv1.CreateUserInfo{
		NickName: user.NickName,
		PassWord: user.PassWord,
		Mobile:   user.Mobile,
	}
	userRsp, err := u.uc.CreateUser(ctx, protoUser)
	if err != nil {
		return err
	}
	user.ID = uint64(userRsp.Id)
	return err
}

func (u *users) Update(ctx context.Context, user *data.User) error {
	protoUser := &upbv1.UpdateUserInfo{
		Id:       int32(user.ID),
		NickName: user.NickName,
		Gender:   user.Gender,
		BirthDay: uint64(user.Birthday.Unix()),
	}
	_, err := u.uc.UpdateUser(ctx, protoUser)
	if err != nil {
		return err
	}
	return nil
}

func (u *users) Get(ctx context.Context, userID uint64) (data.User, error) {
	user, err := u.uc.GetUserById(ctx, &upbv1.IdRequest{Id: int32(userID)})
	if err != nil {
		return data.User{}, err
	}
	return data.User{
		ID:       uint64(user.Id),
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Birthday: itime.Time{time.Unix(int64(user.BirthDay), 0)},
		Gender:   user.Gender,
		Role:     user.Role,
		PassWord: user.PassWord,
	}, nil
}

func (u *users) GetByMobile(ctx context.Context, mobile string) (data.User, error) {
	user, err := u.uc.GetUserByMobile(ctx, &upbv1.MobileRequest{Mobile: mobile})
	if err != nil {
		return data.User{}, err
	}
	return data.User{
		ID:       uint64(user.Id),
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Birthday: itime.Time{time.Unix(int64(user.BirthDay), 0)},
		Gender:   user.Gender,
		Role:     user.Role,
		PassWord: user.PassWord,
	}, nil
}

func (u *users) CheckPassWord(ctx context.Context, password, encryptedPwd string) error {
	cres, err := u.uc.CheckPassWord(ctx, &upbv1.PasswordCheckInfo{
		Password:          password,
		EncryptedPassword: encryptedPwd,
	})
	if err != nil {
		return err
	}
	if cres.Success {
		return nil
	}
	return errors.WithCode(code.ErrUserPasswordIncorrect, "密码错误")
}

var _ data.UserData = &users{}
