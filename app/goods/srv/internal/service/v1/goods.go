package v1

import (
	"context"
	"sync"

	"github.com/zeromicro/go-zero/core/mr"

	proto "mydev/api/goods/v1"
	v1 "mydev/app/goods/srv/internal/data/v1"
	v12 "mydev/app/goods/srv/internal/data_search/v1"
	"mydev/app/goods/srv/internal/domain/do"
	"mydev/app/goods/srv/internal/domain/dto"
	metav1 "mydev/pkg/common/meta/v1"
	"mydev/pkg/log"
)

type GoodsSrv interface {
	// 商品列表
	List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error)

	// 商品详情
	Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error)

	// 创建商品
	Create(ctx context.Context, goods *dto.GoodsDTO) error

	// 更新商品
	Update(ctx context.Context, goods *dto.GoodsDTO) error

	// 删除商品
	Delete(ctx context.Context, ID uint64) error

	//批量查询商品
	BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error)
}

type goodsService struct {
	//工厂
	data       v1.DataFactory
	serachData v12.SearchFactory
	//serachData v12.GoodsStore
	//categoryData v1.CategoryStore
	//brandData    v1.BrandsStore
}

func NewGoodsService(dataFactory v1.DataFactory, serachData v12.SearchFactory) GoodsSrv {
	return &goodsService{
		data:       dataFactory,
		serachData: serachData,
		//categoryData: categoryData,
		//brandData:    brandData,
	}
}
func newGoodsService(srv *service) GoodsSrv {
	return &goodsService{
		data:       srv.data,
		serachData: srv.dataSearch,
		//data:       dataFactory,
		//serachData: serachData,
		//categoryData: categoryData,
		//brandData:    brandData,
	}
}

// 遍历树结构
func retrieveIDs(category *do.CategoryDO) []uint64 {
	ids := []uint64{}
	if category == nil || category.ID == 0 {
		return ids
	}
	ids = append(ids, uint64(category.ID))
	for _, child := range category.SubCategory {
		subids := retrieveIDs(child)
		ids = append(ids, subids...)
	}
	return ids
}

func (gs *goodsService) List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error) {
	searchReq := v12.GoodsFilterRequest{
		GoodsFilterRequest: req,
	}
	//通过父类查找到所有的子分类
	if req.TopCategory > 0 {
		category, err := gs.data.Categorys().Get(ctx, uint64(req.TopCategory))
		if err != nil {
			log.Errorf("categoryData.Get error: %v", err)
			return nil, err
		}
		//转interface
		var ids []interface{}
		for _, value := range retrieveIDs(category) {
			ids = append(ids, value)
		}
		searchReq.CategoryIDs = ids
	}
	//通过es人性化搜索
	goodsList, err := gs.serachData.Goods().Search(ctx, &searchReq)
	if err != nil {
		log.Errorf("searchReq.Search error: %v", err)
		return nil, err
	}
	//只拿id进行后续查询就可以
	goodsIDs := []uint64{}
	for _, value := range goodsList.Items {
		goodsIDs = append(goodsIDs, uint64(value.ID))
	}
	//通过id批量查询mysql数据
	goods, err := gs.data.Goods().ListByIDs(ctx, goodsIDs, orderby)
	if err != nil {
		log.Errorf("ListByIDs error: %v", err)
		return nil, err
	}
	var ret dto.GoodsDTOList
	ret.TotalCount = int(goodsList.TotalCount)
	for _, value := range goods.Items {
		ret.Items = append(ret.Items, &dto.GoodsDTO{
			GoodsDO: *value,
		})
	}
	return &ret, nil

}

func (gs *goodsService) Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error) {
	goods, err := gs.data.Goods().Get(ctx, ID)
	if err != nil {
		log.Errorf("data.Get err:%v", err)
		return nil, err
	}
	return &dto.GoodsDTO{
		GoodsDO: *goods,
	}, nil
}

func (gs *goodsService) Create(ctx context.Context, goods *dto.GoodsDTO) error {
	//数据先写mysql，然后写es
	_, err := gs.data.Brands().Get(ctx, uint64(goods.BrandsID))
	if err != nil {
		return err
	}

	_, err = gs.data.Brands().Get(ctx, uint64(goods.CategoryID))
	if err != nil {
		return err
	}
	/*
		之前的入es的方案是给gorm添加aftercreate
		为了不依赖于gorm中的钩子方法：
		1.分布式事务， 异构数据库的事务， 基于可靠消息最终一致性
			比较重的方案： 每次都要发送一个事务消息
		2.在数据库入库之前启动一个mysql事务，在事务中es入成功才能进行mysql入库，但是data层不能依赖于service，你拿不到Begin。所以 得在data层再声名一个方法是返回db.Begin()
		3.延迟双写：将数据先写入MySQL中，然后异步地写入ES，这样能减轻写入ES的压力，也能降低MySQL事务回滚的概率。但是在数据同步的时间窗口内，MySQL中的数据和ES中的数据可能不一致。
		4.基于日志同步：在MySQL中开启binlog，然后使用binlog来同步数据到ES中。这种方法可以保证MySQL和ES中的数据一致性，但是实现难度较大。
		5.使用ETL工具：使用一些ETL工具，如StreamSets、DataX等，将MySQL中的数据实时同步到ES中。这种方法可以保证数据的实时性和一致性，但是需要使用额外的工具和维护成本。
		6.使用数据库中间件：一些数据库中间件，如MyCat、TDDL等，可以支持将MySQL中的数据同步到ES中。这种方法需要引入额外的中间件，但可以降低开发成本和实现难度。
		选择哪种方法取决于你的具体场景和需求，需要综合考虑实时性、一致性、可靠性、复杂度和成本等因素。
		没有一种方法是完美的，每种方法都有其优缺点和适用场景。选择哪种方法最适合取决于具体的业务场景和需求。
		比如：
		如果对数据的实时性和一致性要求较高，可以考虑使用基于数据库中间件或ETL工具的方案；
		如果对数据的一致性要求很高，可以考虑使用基于日志同步的方案；
		如果对数据的实时性要求不高，但是要求系统高可用和容错能力（即使发生故障或错误，也能够保证系统的正常运行和数据的正确性。）可以考虑使用基于消息的方案。
		这里预留了1、2方案的接口
	*/
	//err = gs.data.Create(ctx, &goods.GoodsDO)
	txn := gs.data.Begin()
	//这个事务要非常小心 commit和rollback的使用
	//如果超时了呢？代码崩掉？网络阻塞？(只是因为网络阻塞 返回了err 并不能证明es插入失败，导致es插入了 mysql没插入)
	//如果程序崩掉了 要捕获recover释放txn
	//这种方案是对我们一致性不高的时候使用的 因为多一条垃圾也无所谓
	defer func() { //很重要
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("goodsService.Create panic: %v", err)
			return
		}
	}()
	//注意传入指针
	err = gs.data.Goods().CreateInTxn(ctx, txn, &goods.GoodsDO)
	if err != nil {
		log.Errorf("data.CreateInTxn err: %v", err)
		txn.Rollback()
		return err
	}

	searchDO := do.GoodsSearchDO{
		ID:          goods.ID,
		CategoryID:  goods.CategoryID,
		BrandsID:    goods.BrandsID,
		OnSale:      goods.OnSale,
		ShipFree:    goods.ShipFree,
		IsNew:       goods.IsNew,
		IsHot:       goods.IsHot,
		Name:        goods.Name,
		ClickNum:    goods.ClickNum,
		SoldNum:     goods.SoldNum,
		FavNum:      goods.FavNum,
		MarketPrice: goods.MarketPrice,
		GoodsBrief:  goods.GoodsBrief,
		ShopPrice:   goods.ShopPrice,
	}
	err = gs.serachData.Goods().Create(ctx, &searchDO)
	if err != nil {
		txn.Rollback()
		return err
	}
	txn.Commit()
	return nil
}

func (gs *goodsService) Update(ctx context.Context, goods *dto.GoodsDTO) error {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsService) Delete(ctx context.Context, ID uint64) error {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsService) BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error) {
	//go-zero 非常好用，但是我们自己去做的并发的话 - 一次性启动多个goroutine
	var ret []*dto.GoodsDTO
	var mu sync.Mutex
	var callFuncs []func() error
	//如果使用map的话要注意他不是线程安全的，使用sync.Map是线程安全的
	//sync.Map.Load()
	//sync.Map.Store()
	for _, value := range ids {
		//大坑
		tem := value
		callFuncs = append(callFuncs, func() error {
			goodsDTO, err := gs.Get(ctx, tem)
			//注意线程安全问题
			mu.Lock()
			ret = append(ret, goodsDTO)
			mu.Unlock()
			return err
		})
	}
	err := mr.Finish(callFuncs...)
	if err != nil {
		return nil, err
	}
	return ret, nil
	//ds, err := gs.data.ListByIDs(ctx, ids, []string{})
	//if err != nil {
	//	return nil, err
	//}
	//var ret []*dto.GoodsDTO
	//for _, value := range ds.Items {
	//	ret = append(ret, &dto.GoodsDTO{
	//		GoodsDO: *value,
	//	})
	//}
}

var _ GoodsSrv = &goodsService{}
