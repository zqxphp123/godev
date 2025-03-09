package v1

import (
	v1 "mydev/app/goods/srv/internal/data/v1"
	v12 "mydev/app/goods/srv/internal/data_search/v1"
)

type ServiceFactory interface {
	Goods() GoodsSrv
	Banner() BannerSrv
}

type service struct {
	data       v1.DataFactory
	dataSearch v12.SearchFactory
}

func (s *service) Goods() GoodsSrv {
	return newGoodsService(s)
}
func (s *service) Banner() BannerSrv {
	return newBannerService(s)
}
func NewService(store v1.DataFactory, dataSearch v12.SearchFactory) *service {
	return &service{data: store, dataSearch: dataSearch}
}

var _ ServiceFactory = &service{}
