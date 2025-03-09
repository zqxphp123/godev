package v1

import (
	"context"

	proto "mydev/api/goods/v1"
	"mydev/app/goods/srv/internal/domain/do"
)

// 请求数据不合适可以做一层封装
type GoodsFilterRequest struct {
	*proto.GoodsFilterRequest
	CategoryIDs []interface{}
}

type GoodsStore interface {
	Create(ctx context.Context, goods *do.GoodsSearchDO) error
	Delete(ctx context.Context, ID uint64) error
	Update(ctx context.Context, goods *do.GoodsSearchDO) error
	Search(ctx context.Context, request *GoodsFilterRequest) (*do.GoodsSearchDOList, error)
}
