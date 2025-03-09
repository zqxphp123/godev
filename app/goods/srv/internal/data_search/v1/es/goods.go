package es

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	v1 "mydev/app/goods/srv/internal/data_search/v1"
	"mydev/app/goods/srv/internal/domain/do"
	"mydev/app/pkg/code"
	"mydev/pkg/errors"
	"strconv"
)

type goods struct {
	esClient *elastic.Client
}

func NewGoodsSearch(esClient *elastic.Client) *goods {
	return &goods{esClient: esClient}
}
func newGoodsSearch(search *dataSearch) *goods {
	return &goods{esClient: search.esClient}
}

func (g *goods) Create(ctx context.Context, goods *do.GoodsSearchDO) error {
	_, err := g.esClient.Index().
		Index(goods.GetIndexName()).
		Id(strconv.Itoa(int(goods.ID))).
		BodyJson(&goods).
		Do(context.TODO())
	return err
}

func (g goods) Delete(ctx context.Context, ID uint64) error {
	_, err := g.esClient.Delete().
		Index(do.GoodsSearchDO{}.GetIndexName()).
		Id(strconv.Itoa(int(ID))).
		Refresh("true").
		Do(ctx)
	return err
}

func (g goods) Update(ctx context.Context, goods *do.GoodsSearchDO) error {
	err := g.Delete(ctx, uint64(goods.ID))
	if err != nil {
		return err
	}
	err = g.Create(ctx, goods)
	if err != nil {
		return err
	}
	return nil
}

func (g goods) Search(ctx context.Context, req *v1.GoodsFilterRequest) (*do.GoodsSearchDOList, error) {
	//match bool 复合查询
	q := elastic.NewBoolQuery()
	if req.KeyWords != "" {
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}

	if req.PriceMin > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}

	if req.Brand > 0 {
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}

	if req.TopCategory > 0 {
		q = q.Filter(elastic.NewTermsQuery("category_id", req.CategoryIDs...))
	}
	//分页
	if req.Pages == 0 {
		req.Pages = 1
	}

	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}
	res, err := g.esClient.Search().Index(do.GoodsSearchDO{}.GetIndexName()).
		Query(q).
		From(int(req.Pages-1) * int(req.PagePerNums)).
		Size(int(req.PagePerNums)).
		Do(ctx)

	var ret do.GoodsSearchDOList

	ret.TotalCount = res.Hits.TotalHits.Value
	for _, value := range res.Hits.Hits {
		goodse := do.GoodsSearchDO{}
		err = json.Unmarshal(value.Source, &goodse)
		if err != nil {
			return nil, errors.WithCode(code.ErrEsUnmarshal, err.Error())
		}
		ret.Items = append(ret.Items, &goodse)
	}

	return &ret, nil
}

var _ v1.GoodsStore = &goods{}
