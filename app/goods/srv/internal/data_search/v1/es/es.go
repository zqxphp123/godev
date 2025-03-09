package es

import (
	"github.com/olivere/elastic/v7"
	v1 "mydev/app/goods/srv/internal/data_search/v1"
	"mydev/app/pkg/options"
	"mydev/pkg/db"
	"mydev/pkg/errors"
	"sync"
)

var (
	//esClient      *elastic.Client
	searchFactory v1.SearchFactory
	once          sync.Once
)

type dataSearch struct {
	esClient *elastic.Client
}

func (ds *dataSearch) Goods() v1.GoodsStore {
	return newGoodsSearch(ds)
}

var _ v1.SearchFactory = &dataSearch{}

func GetSearchFactoryOr(opts *options.EsOptions) (v1.SearchFactory, error) {
	if opts == nil && searchFactory == nil {
		return nil, errors.New("failed to get es client")
	}
	//var err error
	once.Do(func() {
		esOpt := db.EsOptions{
			Host: opts.Host,
			Port: opts.Port,
		}
		esClient, err := db.NewEsClient(&esOpt)
		if err != nil {
			return
		}
		searchFactory = &dataSearch{esClient: esClient}
	})
	if searchFactory == nil {
		return nil, errors.New("failed to get ES client")
	}
	return searchFactory, nil
}
