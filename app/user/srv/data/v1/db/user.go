package db

import (
	"context"

	"gorm.io/gorm"

	"mydev/app/pkg/code"
	dv1 "mydev/app/user/srv/data/v1"
	code2 "mydev/gmicro/code"
	metav1 "mydev/pkg/common/meta/v1"
	"mydev/pkg/errors"
)

type users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *users {
	return &users{db: db}
}

// List
//
//	@Description: 获取用户列表，凡是列表页返回的时候都应该返回总共有多少个
//	@receiver u
//	@param ctx
//	@param opts: 分页
//	@return *dv1.UserDOList
//	@return error
func (u *users) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*dv1.UserDOList, error) {
	ret := dv1.UserDOList{}
	//分页
	var limit, offset int
	if opts.PageSize == 0 {
		limit = 10
	} else {
		limit = opts.PageSize
	}
	if opts.Page > 0 {
		offset = (opts.Page - 1) * limit
	}
	//排序
	query := u.db
	for _, value := range orderby {
		//坑：如果db改掉了？
		//u.db=u.db.Order(value)
		query = query.Order(value)
	}
	//查询 - 发起多个请求
	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return &ret, nil
}

// GetByMobile
//
//	@Description: 通过手机号码查询用户
//	@receiver u
//	@param ctx
//	@param mobile: 用户手机号
//	@return *dv1.UserDO
//	@return error
func (u *users) GetByMobile(ctx context.Context, mobile string) (*dv1.UserDO, error) {
	user := dv1.UserDO{}

	//err是gorm的error这种error我们尽量不要抛出去
	err := u.db.Where("mobile = ?", mobile).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

// GetByID
//
//	@Description: 通过用户ID查询用户
//	@receiver u
//	@param ctx
//	@param id: 用户id
//	@return *dv1.UserDO
//	@return error
func (u *users) GetByID(ctx context.Context, id uint64) (*dv1.UserDO, error) {
	user := dv1.UserDO{}

	//err是gorm的error这种error我们尽量不要抛出去
	err := u.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

// Create
//
//	@Description: 创建用户
//	@receiver u
//	@param ctx
//	@param user: 用户DO
//	@return error
func (u *users) Create(ctx context.Context, user *dv1.UserDO) error {
	tx := u.db.Create(user)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

// Update
//
//	@Description: 更新用户
//	@receiver u
//	@param ctx
//	@param user: 用户DO
//	@return error
func (u *users) Update(ctx context.Context, user *dv1.UserDO) error {
	tx := u.db.Save(user)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func newUsers(db *gorm.DB) *users {
	return &users{db: db}
}

// 校验  函数签名
var _ dv1.UserStore = &users{}
