package v1

import (
	"context"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	v1 "mydev/app/inventory/srv/internal/data/v1"
	"mydev/app/inventory/srv/internal/domain/do"
	"mydev/app/inventory/srv/internal/domain/dto"
	"mydev/app/pkg/code"
	"mydev/app/pkg/options"
	"mydev/pkg/errors"
	"mydev/pkg/log"
	"sort"
)

const (
	inventoryLockPrefix = "inventory_"
	orderLockProfix     = "order_"
)

type InventorySrv interface {
	//设置库存
	Create(ctx context.Context, inv *dto.InventoryDTO) error
	//根据商品id查询库存
	Get(ctx context.Context, goodsID uint64) (*dto.InventoryDTO, error)
	//扣减库存
	Sell(ctx context.Context, ordersn string, detail []do.GoodsDetail) error
	//归还库存
	Reback(ctx context.Context, ordersn string, detail []do.GoodsDetail) error
}
type inventoryService struct {
	data      v1.DataFactory
	redisOpts *options.RedisOptions
	pool      redsyncredis.Pool
}

func (is *inventoryService) Create(ctx context.Context, inv *dto.InventoryDTO) error {
	return is.data.Inventorys().Create(ctx, &inv.InventoryDO)
}

func (is *inventoryService) Get(ctx context.Context, goodsID uint64) (*dto.InventoryDTO, error) {
	inv, err := is.data.Inventorys().Get(ctx, goodsID)
	if err != nil {
		return nil, err
	}
	return &dto.InventoryDTO{*inv}, nil
}

func (is *inventoryService) Sell(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单%s扣减库存", ordersn)
	//数据一致性要求极高，尽可能用多一点的技术完成
	rs := redsync.New(is.pool)
	/*
		实际上批量扣减库存的时候，我们经常会先按照商品的id排序，然后从小到大逐个扣减库存，这样可以减少锁的竞争，
		例如：如果线程1要对商品A进行扣减操作，线程2要对商品B进行扣减操作，那么线程1会先获取商品A的锁，然后扣减库存，最后释放锁；而线程2会等待线程1释放商品A的锁后，才能获取商品B的锁，进行扣减操作。
		这样，不同线程对不同商品进行扣减操作，就避免了锁的竞争问题。如果不排序1扣减a b c ；2扣减c b a，会造成同时对商品a进行枪锁，从而引发锁竞争，降低系统的并发性能。
		还能解决一下问题：
		避免死锁：如果多个线程同时对多个商品进行扣减操作，而没有排序，则可能会引发死锁问题。例如，线程1占用了商品A的锁，等待商品B的锁，而线程2占用了商品B的锁，等待商品A的锁，这样就出现了死锁情况。而通过按照商品ID排序，可以保证所有线程都按照相同的顺序获取商品的锁，避免了死锁问题的发生。
		确保数据一致性：如果没有排序而直接进行库存扣减操作，可能会导致数据不一致的问题。例如，线程1对商品A进行了库存扣减操作，而线程2也对商品A进行了扣减操作，但是由于线程1和线程2操作的顺序不确定，可能会导致库存数据不一致。而通过对商品ID进行排序，可以保证对于同一个商品的扣减操作按照相同的顺序进行，避免了数据不一致的问题。
	*/
	var detail = do.GoodsDetailList(details)
	//实现len i<j   i,j=j,i   即可
	sort.Sort(detail)
	txn := is.data.Begin()
	//使用事务必须要捕获异常
	defer func() {
		if err := recover(); err != nil {
			_ = txn.Rollback()
			log.Error("事务进行中出现异常，回滚")
			return
		}
	}()
	sellDetail := do.StockSellDetailDO{
		OrderSn: ordersn,
		Status:  1,
		Detail:  detail,
	}
	for _, goodsInfo := range detail {
		mutex := rs.NewMutex(inventoryLockPrefix + ordersn)
		if err := mutex.Lock(); err != nil {
			log.Errorf("订单%s获取锁失败", ordersn)
			return err
		}
		//判断库存是否存在
		inv, err := is.data.Inventorys().Get(ctx, uint64(goodsInfo.Goods))
		if err != nil {
			log.Errorf("订单%s获取库存失败", ordersn)
			return err
		}
		//判断库存是否充足
		if inv.Stocks < goodsInfo.Num {
			txn.Rollback() //回滚
			log.Errorf("商品%d库存%d不足, 现有库存：%d", goodsInfo.Goods, goodsInfo.Num, inv.Stocks)
			return errors.WithCode(code.ErrInvNotEnough, "库存不足")
		}
		inv.Stocks -= goodsInfo.Num
		err = is.data.Inventorys().Reduce(ctx, txn, uint64(goodsInfo.Goods), int(goodsInfo.Num))
		if err != nil {
			txn.Rollback() //回滚
			log.Errorf("订单%s扣减库存失败", ordersn)
			return err
		}
		//释放锁
		if _, err = mutex.Unlock(); err != nil {
			txn.Rollback() //回滚
			log.Errorf("订单%s释放锁出现异常", ordersn)
		}
	}
	err := is.data.Inventorys().CreateStockSellDetail(ctx, txn, &sellDetail)
	if err != nil {
		txn.Rollback() //回滚
		log.Errorf("订单%s创建扣减库存记录失败", ordersn)
		return err
	}
	txn.Commit()
	return nil
}

func (is *inventoryService) Reback(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单%s库存归还", ordersn)
	rs := redsync.New(is.pool)
	txn := is.data.Begin()
	//使用事务必须要捕获异常
	defer func() {
		if err := recover(); err != nil {
			_ = txn.Rollback()
			log.Error("事务进行中出现异常，回滚")
			return
		}
	}()
	//库存归还的时候有不少细节
	//1.主动取消 2.网络问题引起的重试 3.超时取消 4.退款取消 5.系统故障取消 6.风险控制取消
	//加锁
	mutex := rs.NewMutex(orderLockProfix + ordersn)
	if err := mutex.Lock(); err != nil {
		txn.Rollback() //回滚
		log.Errorf("订单%s获取锁失败", ordersn)
		return err
	}
	sellDetail, err := is.data.Inventorys().GetSellDetail(ctx, txn, ordersn)
	if err != nil {
		txn.Rollback() //回滚锁
		_, err = mutex.Unlock()
		if err != nil {
			log.Errorf("订单%s释放锁出现异常", ordersn)
			return err
		}
		if errors.IsCode(err, code.ErrInvSellDetailNotFound) {
			log.Errorf("订单%s扣减库记录不存在,忽略", ordersn)
			return nil
		}
		log.Errorf("订单%s获取扣减库存记录失败", ordersn)
		return err
	}
	if sellDetail.Status == 2 {
		log.Infof("订单%s扣减库存记录已经归还，忽略", ordersn)
	}
	var detail = do.GoodsDetailList(details)
	sort.Sort(detail)
	for _, goodsInfo := range detail {
		inv, err := is.data.Inventorys().Get(ctx, uint64(goodsInfo.Goods))
		if err != nil {
			txn.Rollback() //回滚
			log.Errorf("订单%s获取库存失败", ordersn)
			return err
		}
		inv.Stocks += goodsInfo.Num
		err = is.data.Inventorys().Increase(ctx, txn, uint64(goodsInfo.Goods), int(goodsInfo.Num))
		if err != nil {
			txn.Rollback() //回滚
			log.Errorf("订单%s归还库存失败", ordersn)
			return err
		}
	}
	err = is.data.Inventorys().UpdateStockSellDetailStatus(ctx, txn, ordersn, 2)
	if err != nil {
		txn.Rollback() //回滚
		log.Errorf("订单%s更新扣减库存记录失败", ordersn)
		return err
	}
	txn.Commit()
	return nil
}

func newInventoryService(s *service) *inventoryService {
	return &inventoryService{data: s.data, redisOpts: s.redisOpts, pool: s.pool}
}

var _ InventorySrv = &inventoryService{}
