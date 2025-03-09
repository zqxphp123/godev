package db

import (
	"context"

	"gorm.io/gorm"

	v1 "mydev/app/goods/srv/internal/data/v1"
	"mydev/app/goods/srv/internal/domain/do"
	"mydev/app/pkg/code"
	code2 "mydev/gmicro/code"
	metav1 "mydev/pkg/common/meta/v1"
	"mydev/pkg/errors"
)

type goods struct {
	db *gorm.DB
}

// 向外提供
func NewGoods(db *gorm.DB) *goods {
	return &goods{
		db: db,
	}
}

// 提供工厂
func newGoods(factory *mysqlFactory) *goods {
	return &goods{
		db: factory.db,
	}
}
func (g *goods) Get(ctx context.Context, ID uint64) (*do.GoodsDO, error) {
	good := &do.GoodsDO{}

	err := g.db.Preload("Brands").Preload("Category").Where(ID).First(&good).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrGoodsNotFound, err.Error())
		}
	}
	return good, nil
}

func (g *goods) ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error) {
	ret := &do.GoodsDOList{}
	//排序
	query := g.db.Preload("Category").Preload("Brands")
	for _, value := range orderby {
		//坑：如果db改掉了？
		//u.db=u.db.Order(value)
		query = query.Order(value)
	}
	//查询 - 发起多个请求
	d := query.Where("id in ?", ids).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

func (g *goods) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.GoodsDOList, error) {
	ret := &do.GoodsDOList{}
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
	query := g.db.Preload("Category").Preload("Brands")
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
	return ret, nil
}

func (g *goods) Create(ctx context.Context, goods *do.GoodsDO) error {
	tx := g.db.Create(goods)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) CreateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	tx := txn.Create(goods)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Update(ctx context.Context, goods *do.GoodsDO) error {
	tx := g.db.Save(goods)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) UpdateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	tx := txn.Save(goods)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Delete(ctx context.Context, ID uint64) error {
	return g.db.Where("id = ?", ID).Delete(&do.GoodsDO{}).Error
}

func (g *goods) DeleteInTxn(ctx context.Context, txn *gorm.DB, ID uint64) error {
	return txn.Where("id = ?", ID).Delete(&do.GoodsDO{}).Error
}

func (g *goods) Begin() *gorm.DB {
	return g.db.Begin()
}

var _ v1.GoodsStore = &goods{}
