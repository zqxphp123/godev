package db

import (
	"context"
	code2 "mydev/gmicro/code"
	"mydev/pkg/errors"

	"gorm.io/gorm"

	v1 "mydev/app/goods/srv/internal/data/v1"
	"mydev/app/goods/srv/internal/domain/do"
	metav1 "mydev/pkg/common/meta/v1"
)

type banners struct {
	db *gorm.DB
}

func NewBanner(db *gorm.DB) *banners {
	return &banners{
		db: db,
	}
}
func newBanner(factory *mysqlFactory) *banners {
	return &banners{
		db: factory.db,
	}
}
func (b *banners) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BannerList, error) {
	listBaner := &do.BannerList{}

	baners := []do.BannerDO{}
	dbBaner := b.db.Find(&baners)
	listBaner.TotalCount = dbBaner.RowsAffected
	for _, banner := range baners {
		listBaner.Items = append(listBaner.Items, &banner)
	}
	return listBaner, nil
}

func (b *banners) Create(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error {
	db := b.db.Create(banner)
	if db.Error != nil {
		return errors.WithCode(code2.ErrDatabase, db.Error.Error())
	}
	return nil
}

func (b *banners) Update(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error {
	db := b.db.Save(banner)
	if db.Error != nil {
		return errors.WithCode(code2.ErrDatabase, db.Error.Error())
	}
	return nil
}

func (b *banners) Delete(ctx context.Context, ID uint64) error {
	return b.db.Where("id = ?", ID).Delete(&do.BannerDO{}).Error
}

var _ v1.BannerStore = &banners{}
