package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mydev/app/pkg/code"
	"mydev/pkg/errors"
	"sync"

	"mydev/app/pkg/options"
)

var (
	dbFactory *gorm.DB
	once      sync.Once
)

// 这个方法会返回gorm连接
// 还不够
// 这个方法应该返回的是全局的一个变量，如果一开始的时候没有初始化好，那么就初始化一次，后续呢直接拿到这个变量
// 单例模式  演进成 sync.Once
func GetDBFactoryOr(mysqlOpts *options.MySQLOptions) (*gorm.DB, error) {
	if mysqlOpts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlOpts.Username,
			mysqlOpts.Password,
			mysqlOpts.Host,
			mysqlOpts.Port,
			mysqlOpts.Database)
		dbFactory, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}
		sqlDB, _ := dbFactory.DB()
		//允许连接多少个mysql
		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		//允许最大的空闲的连接数
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		//重用连接的最大时长
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifetime)

	})
	if dbFactory == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}
	return dbFactory, nil

}
