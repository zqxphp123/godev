package service

import (
	"context"
	"github.com/dtm-labs/client/dtmgrpc"
	proto2 "mydev/api/goods/v1"
	proto "mydev/api/inventory/v1"
	proto3 "mydev/api/order/v1"
	v12 "mydev/app/order/srv/internal/data/v1"
	"mydev/app/order/srv/internal/domain/do"
	"mydev/app/order/srv/internal/domain/dto"
	"mydev/app/pkg/code"
	"mydev/app/pkg/options"
	v1 "mydev/pkg/common/meta/v1"
	"mydev/pkg/errors"
	"mydev/pkg/log"
)

const (
	goodsBusi = "discovery:///mydev-order-srv"
	invBusi   = "discovery:///mydev-inventory-srv"
)

type OrderSrv interface {
	Get(ctx context.Context, orderSn string) (*dto.OrderDTO, error)
	List(ctx context.Context, userID uint64, meta v1.ListMeta, orderby []string) (*dto.OrderDTOList, error)
	Submit(ctx context.Context, order *dto.OrderDTO) error
	Create(ctx context.Context, order *dto.OrderDTO) error
	CreateCom(ctx context.Context, order *dto.OrderDTO) error //这是create的补偿
	Update(ctx context.Context, order *dto.OrderDTO) error
}
type orderService struct {
	data    v12.DataFactory
	dtmOpts *options.DtmOptions
}

func (os *orderService) Get(ctx context.Context, orderSn string) (*dto.OrderDTO, error) {
	order, err := os.data.Orders().Get(ctx, orderSn)
	if err != nil {
		return nil, err
	}
	return &dto.OrderDTO{*order}, nil
}

func (os *orderService) List(ctx context.Context, userID uint64, meta v1.ListMeta, orderby []string) (*dto.OrderDTOList, error) {
	orders, err := os.data.Orders().List(ctx, userID, meta, orderby)
	if err != nil {
		return nil, err
	}
	var ret = dto.OrderDTOList{}
	for _, value := range orders.Items {
		ret.Items = append(ret.Items, &dto.OrderDTO{*value})
	}
	return &ret, nil
}

func (os *orderService) Submit(ctx context.Context, order *dto.OrderDTO) error {
	//先从购物车中获取商品信息 - 填充数据
	list, err := os.data.ShopCarts().List(ctx, uint64(order.User), true, v1.ListMeta{}, []string{})
	if err != nil {
		log.Errorf("获取购物车信息失败，err:%s", err)
		return err
	}
	if len(list.Items) == 0 {
		log.Error("购物车中没有商品")
		return errors.WithCode(code.ErrNoGoodsSelect, "购物车中没有商品")
	}
	var orderGoods []*do.OrderGoods
	var orderItems []*proto3.OrderItemResponse
	for _, value := range list.Items {
		orderGoods = append(orderGoods, &do.OrderGoods{
			Goods: value.Goods,
			Nums:  value.Nums,
		})
		orderItems = append(orderItems, &proto3.OrderItemResponse{
			GoodsId: value.Goods,
			Nums:    value.Nums,
		})
	}
	order.OrderGoods = orderGoods

	//基于可靠消息最终一致性方法，saga事务来解决订单生成的问题

	var goodsInfo []*proto.GoodsInvInfo
	for _, value := range order.OrderGoods {
		goodsInfo = append(goodsInfo, &proto.GoodsInvInfo{
			GoodsId: value.Goods,
			Num:     value.Nums,
		})
	}
	//orderSn := generateOrderSn(order.User)
	req := &proto.SellInfo{
		GoodsInfo: goodsInfo,
		OrderSn:   order.OrderSn,
	}

	oReq := &proto3.OrderRequest{
		UserId:     order.User,
		Address:    order.Address,
		Name:       order.SignerName,
		Mobile:     order.SingerMobile,
		Post:       order.Post,
		OrderSn:    order.OrderSn,
		OrderItems: orderItems,
	}

	saga := dtmgrpc.NewSagaGrpc(os.dtmOpts.GrpcServer, order.OrderSn).
		Add(invBusi+"/Inventory/Sell", invBusi+"/Inventory/Reback", req).
		Add(goodsBusi+"/Order/CreateOrder", goodsBusi+"/Order/CreateOrderCom", oReq)
	saga.WaitResult = true

	err = saga.Submit()
	//通过gid(OrderSn)查询一下，当前的状态如何状态一直值Submitted那么就你一直不要给前端返回,如果是failed那么你提示给前端说下单失败，重新下单
	return err
}

func (os *orderService) Create(ctx context.Context, order *dto.OrderDTO) error {
	/*
		1.生成orderinfo表
		2.生成ordergoods表
		3.根据order找到对应的购物车条目 进行删除。这一步在本地事务
	*/
	var goodsids []int32
	for _, value := range order.OrderGoods {
		goodsids = append(goodsids, value.Goods)
	}
	//return status.Error(codes.Aborted, "create order failed")
	//debug过程或dtm中的退避算计中可能会触发context的超时机制
	goods, err := os.data.Goods().BatchGetGoods(context.Background(), &proto2.BatchGoodsIdInfo{Id: goodsids})
	if err != nil {
		log.Errorf("批量获取商品信息失败,goodsids: %v,err: %v", goodsids, err.Error())
		return err //这个不是abort  也就是说  会不停的重试
	}
	if len(goods.Data) != len(goodsids) {
		log.Errorf("批量获取商品信息失败,goodsids: %v,返回值: %v,err: %v", goodsids, goods.Data, err.Error())
		return errors.WithCode(code.ErrGoodsNotFound, "商品不存在或部分不存在")
	}
	var goodsMap = make(map[int32]*proto2.GoodsInfoResponse)
	for _, value := range goods.Data {
		goodsMap[value.Id] = value
	}
	var orderAmount float32
	for _, value := range order.OrderGoods {
		orderAmount += goodsMap[value.Goods].ShopPrice * float32(value.Nums)
		value.GoodsName = goodsMap[value.Goods].Name
		value.GoodsPrice = goodsMap[value.Goods].ShopPrice
		value.GoodsImage = goodsMap[value.Goods].GoodsFrontImage
	}
	order.OrderMount = orderAmount
	txn := os.data.Begin()
	defer func() {
		if err := recover(); err != nil {
			_ = txn.Rollback()
			log.Error("新建订单事务进行中出现异常，回滚")
			return
		}
	}()
	err = os.data.Orders().Create(ctx, txn, &order.OrderInfoDO)
	if err != nil {
		txn.Rollback()
		log.Errorf("创建订单失败，err: %v", err.Error())
		return err
	}
	//如果客户在2个客户端 一边提交订单  一边删除购物车怎么办？   可以忽略err
	err = os.data.ShopCarts().DeleteByGoodsIDs(ctx, txn, uint64(order.User), goodsids)
	if err != nil {
		txn.Rollback()
		log.Errorf("删除购物车失败,goodsids: %v,err: %v", goodsids, err.Error())
		return err
	}
	txn.Commit()
	//这里不应该有逻辑 上述应该都成功或失败
	return nil
}

func (os *orderService) CreateCom(ctx context.Context, order *dto.OrderDTO) error {
	/*
		1.删除orderinfo表
		2.删除ordergoods表
		是否需要？3.根据order找到对应的购物车条目 进行补偿。这一步在本地事务
	*/
	//其实不用回滚
	//应该先查询订单时候存在，如果已经存在删除相关记录即可，同时删除购物车记录
	return nil
}

func (os *orderService) Update(ctx context.Context, order *dto.OrderDTO) error {
	//TODO implement me
	panic("implement me")
}

func newOrderService(sv *service) *orderService {
	return &orderService{
		data:    sv.data,
		dtmOpts: sv.dtmOpts,
	}
}

var _ OrderSrv = &orderService{}
