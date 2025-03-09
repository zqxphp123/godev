package db

import (
	"context"
	"gorm.io/gorm"
	v1 "mydev/app/goods/srv/internal/data/v1"
	"mydev/app/goods/srv/internal/domain/do"
	metav1 "mydev/pkg/common/meta/v1"
)

type categoryBrands struct {
	db *gorm.DB
}

func (c categoryBrands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandList, error) {
	//TODO implement me
	panic("implement me")
}

func (c categoryBrands) Create(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error {
	//TODO implement me
	panic("implement me")
}

func (c categoryBrands) Update(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error {
	//TODO implement me
	panic("implement me")
}

func (c categoryBrands) Delete(ctx context.Context, ID uint64) error {
	//TODO implement me
	panic("implement me")
}

var _ v1.GoodsCategoryBrandStore = &categoryBrands{}
