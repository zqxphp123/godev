package v1

import (
	"context"
	"fmt"
	"mydev/app/pkg/code"
	"mydev/pkg/errors"
	"mydev/pkg/log"
	"mydev/pkg/storage"
	"time"

	"github.com/dgrijalva/jwt-go"
	"mydev/app/mydev/api/internal/data"
	"mydev/app/pkg/options"
	"mydev/gmicro/server/restserver/middlewares"
)

type UserDTO struct {
	data.User

	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}
type UserSrv interface {
	MobileLogin(ctx context.Context, mobile, password string) (*UserDTO, error)
	Register(ctx context.Context, mobile, password string, codes string) (*UserDTO, error)
	Updata(ctx context.Context, userDTO *UserDTO) error
	Get(ctx context.Context, userID uint64) (*UserDTO, error)
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
	CheckPassWord(ctx context.Context, password, EncryptedPassword string) (bool, error)
}
type userService struct {
	//ud data.UserData
	data data.DataFactory

	jwtOpts *options.JwtOptions
}

func NewUser(data data.DataFactory, jwtOpts *options.JwtOptions) UserSrv {
	return &userService{data: data, jwtOpts: jwtOpts}
}
func (us *userService) MobileLogin(ctx context.Context, mobile, password string) (*UserDTO, error) {
	//TODO 检查验证码是否正确

	user, err := us.data.Users().GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	//检查密码是否正确
	err = us.data.Users().CheckPassWord(ctx, password, user.PassWord)
	if err != nil {
		return nil, err
	}

	//生成token
	j := middlewares.NewJWT(us.jwtOpts.Key)
	claims := middlewares.CustomClaims{
		//key需要和jwt的key一致
		//这里json能起作用是因为在j.CreateToken里NewWithClaims的时候底层claims会解析成json
		ID:          uint(user.ID),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                                   //签名的生效时间
			ExpiresAt: (time.Now().Local().Add(us.jwtOpts.Timeout)).Unix(), //*天过期
			Issuer:    us.jwtOpts.Realm,
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, err
	}
	return &UserDTO{
		User:      user,
		Token:     token,
		ExpiresAt: (time.Now().Local().Add(us.jwtOpts.Timeout)).Unix(),
	}, nil
}

func (u *userService) Register(ctx context.Context, mobile, password string, codes string) (*UserDTO, error) {
	//验证码校验
	rstore := storage.RedisCluster{}
	value, err := rstore.GetKey(ctx, fmt.Sprintf("%s_%d", mobile, 1))
	if err != nil {
		return nil, errors.WithCode(code.ErrCodeNotExist, "验证码不存在")
	}
	if value != codes {
		return nil, errors.WithCode(code.ErrCodeInCorrect, "验证码不匹配")
	}
	var userDO = &data.User{
		Mobile:   mobile,
		PassWord: password,
	}
	err = u.data.Users().Create(ctx, userDO)
	if err != nil {
		log.Errorf("user register failed: %v", err)
		return nil, err
	}
	//生成token
	j := middlewares.NewJWT(u.jwtOpts.Key)
	claims := middlewares.CustomClaims{
		ID:          uint(userDO.ID),
		NickName:    userDO.NickName,
		AuthorityId: uint(userDO.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                                  //签名的生效时间
			ExpiresAt: (time.Now().Local().Add(u.jwtOpts.Timeout)).Unix(), //*天过期
			Issuer:    u.jwtOpts.Realm,
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, err
	}
	return &UserDTO{
		User:      *userDO,
		Token:     token,
		ExpiresAt: (time.Now().Local().Add(u.jwtOpts.Timeout)).Unix(),
	}, nil
}

func (u *userService) Updata(ctx context.Context, userDTO *UserDTO) error {
	user := &data.User{
		Mobile:   userDTO.Mobile,
		NickName: userDTO.NickName,
		Birthday: userDTO.Birthday,
		Gender:   userDTO.Gender,
	}
	err := u.data.Users().Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) Get(ctx context.Context, userID uint64) (*UserDTO, error) {
	userDO, err := u.data.Users().Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &UserDTO{User: userDO}, nil
}

func (u *userService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) CheckPassWord(ctx context.Context, password, EncryptedPassword string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

var _ UserSrv = &userService{}
