package db

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	v1 "mydev/app/goods/srv/internal/data/v1"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"mydev/app/pkg/code"
	"mydev/app/pkg/options"
	"mydev/pkg/errors"
)

var (
	dbFactory v1.DataFactory
	once      sync.Once
)

// 由工厂持有它 一律注入工厂 不注入Store了 一律由工厂决定
type mysqlFactory struct {
	db *gorm.DB
}

func (mf *mysqlFactory) Begin() *gorm.DB {
	return mf.db.Begin()
}

func (mf *mysqlFactory) Goods() v1.GoodsStore {
	return newGoods(mf)
}

func (mf *mysqlFactory) Categorys() v1.CategoryStore {
	return newCategorys(mf)
}

func (mf *mysqlFactory) Brands() v1.BrandsStore {
	return newBrands(mf)
}

func (mf *mysqlFactory) Banner() v1.BannerStore {
	return newBanner(mf)
}

func (mf *mysqlFactory) CategoryBrands() v1.GoodsCategoryBrandStore {
	//TODO implement me
	panic("implement me")
}

var _ v1.DataFactory = &mysqlFactory{}

// 这个方法会返回gorm连接
// 还不够
// 这个方法应该返回的是全局的一个变量，如果一开始的时候没有初始化好，那么就初始化一次，后续呢直接拿到这个变量
// 单例模式  演进成 sync.Once
func GetDBFactoryOr(mysqlOpts *options.MySQLOptions) (v1.DataFactory, error) {
	if mysqlOpts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlOpts.Username,
			mysqlOpts.Password,
			mysqlOpts.Host,
			mysqlOpts.Port,
			mysqlOpts.Database)
		//可以自己封装log
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				//慢查询阈值可以自己设置
				SlowThreshold:             time.Second,                         // Slow SQL threshold
				LogLevel:                  logger.LogLevel(mysqlOpts.LogLevel), // Log level
				IgnoreRecordNotFoundError: true,                                // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,                                // Don't include params in the SQL log
				Colorful:                  false,                               // Disable color
			},
		)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return
		}
		sqlDB, _ := db.DB()
		dbFactory = &mysqlFactory{
			db: db,
		}
		//允许连接多少个mysql
		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		//允许最大的空闲的连接数
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		//重用连接的最大时长
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifetime)

	})
	var err error
	if dbFactory == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}
	return dbFactory, nil
}
