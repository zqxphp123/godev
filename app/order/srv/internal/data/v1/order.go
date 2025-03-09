package v1

import (
	"context"
	"gorm.io/gorm"
	"mydev/app/order/srv/internal/domain/do"
	v1 "mydev/pkg/common/meta/v1"
)

//	type OrderFilter struct {
//		userID
//		startTime
//		endTime
//	}
type OrderStore interface {
	//订单详情
	Get(ctx context.Context, orderSn string) (*do.OrderInfoDO, error)
	//查询订单列表
	List(ctx context.Context, userID uint64, meta v1.ListMeta, orderby []string) (*do.OrderInfoDOList, error)
	//新建订单
	Create(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error
	//更新订单
	Update(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error
}
