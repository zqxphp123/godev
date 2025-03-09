package v1

import (
	"context"
	gpb "mydev/api/goods/v1"
	"mydev/app/mydev/api/internal/data"
)

type GoodsSrv interface {
	List(ctx context.Context, request *gpb.GoodsFilterRequest) (*gpb.GoodsListResponse, error)
}

type goodsService struct {
	data data.DataFactory
}

func (gs *goodsService) List(ctx context.Context, request *gpb.GoodsFilterRequest) (*gpb.GoodsListResponse, error) {
	rsp, err := gs.data.Goods().GoodsList(ctx, request)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// 这里不用sevice是为了防止循环引用，其实也可以把service放到同一个文件夹下，但是那样文件多了就分不清了  不好看
func NewGoods(data data.DataFactory) *goodsService {
	return &goodsService{data: data}
}

var _ GoodsSrv = &goodsService{}
