package order

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "mydev/api/order/v1"
	"mydev/app/order/srv/internal/domain/do"
	"mydev/app/order/srv/internal/domain/dto"
	"mydev/app/order/srv/internal/service/v1"
	"mydev/pkg/log"
)

type orderServer struct {
	srv service.ServiceFactory
	pb.UnimplementedOrderServer
}

func NewOrderServer(srv service.ServiceFactory) *orderServer {
	return &orderServer{srv: srv}
}

//api层做
//func generateOrderSn(userId int32) string {
//	//订单号的生成规则
//	/*
//		年月日时分秒+用户id+2位随机数
//	*/
//	now := time.Now()
//	rand.New(rand.NewSource(time.Now().UnixNano()))
//	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
//		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
//		userId, rand.Intn(90)+10,
//	)
//	return orderSn
//}

func (os *orderServer) CartItemList(ctx context.Context, info *pb.UserInfo) (*pb.CartItemListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *orderServer) CreateCartItem(ctx context.Context, request *pb.CartItemRequest) (*pb.ShopCartInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *orderServer) UpdateCartItem(ctx context.Context, request *pb.CartItemRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (os *orderServer) DeleteCartItem(ctx context.Context, request *pb.CartItemRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

// 给分布式事务saga调用的，目前没有为api提供的目的
func (os *orderServer) CreateOrder(ctx context.Context, request *pb.OrderRequest) (*emptypb.Empty, error) {
	orderGoods := make([]*do.OrderGoods, len(request.OrderItems))
	for i, item := range request.OrderItems {
		orderGoods[i] = &do.OrderGoods{
			Goods: item.GoodsId,
			Nums:  item.Nums,
		}
	}
	err := os.srv.Orders().Create(ctx, &dto.OrderDTO{
		OrderInfoDO: do.OrderInfoDO{
			User:         request.UserId,
			Address:      request.Address,
			SignerName:   request.Name,
			SingerMobile: request.Mobile,
			Post:         request.Post,
			OrderSn:      request.OrderSn,
			OrderGoods:   orderGoods,
		},
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (os *orderServer) CreateOrderCom(ctx context.Context, request *pb.OrderRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

/*
订单提交的时候应该是先生成订单号
订单号会单独做一个接口，订单查询，以及一系列的关联我们应该采用order_sn，不要再去采用id去关联
*/
func (os *orderServer) SubmitOrder(ctx context.Context, request *pb.OrderRequest) (*emptypb.Empty, error) {
	//从购物车中得到选中的商品 放到service层做

	orderDTO := dto.OrderDTO{
		OrderInfoDO: do.OrderInfoDO{
			User:         request.UserId,
			Address:      request.Address,
			SignerName:   request.Name,
			SingerMobile: request.Mobile,
			Post:         request.Post,
			OrderSn:      request.OrderSn,
		},
	}
	err := os.srv.Orders().Submit(ctx, &orderDTO)
	if err != nil {
		//if errors.IsCode(err, code.ErrNoGoodsSelect) {
		//	return nil, errors.WithCode(code.ErrNoGoodsSelect, "没有选中的商品")
		//}
		log.Errorf("新建订单失败: %s", err.Error())
		return nil, err
	}

	//OrderInfoRsp := &pb.OrderInfoResponse{
	//	UserId:  request.UserId,
	//	Address: request.Address,
	//	Name:    request.Name,
	//	Mobile:  request.Mobile,
	//	Post:    request.Post,
	//	OrderSn: generateOrderSn(request.UserId),
	//}
	return &emptypb.Empty{}, nil
}

func (os *orderServer) OrderList(ctx context.Context, request *pb.OrderFilterRequest) (*pb.OrderListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *orderServer) OrderDetail(ctx context.Context, request *pb.OrderRequest) (*pb.OrderInfoDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *orderServer) UpdateOrderStatus(ctx context.Context, status *pb.OrderStatus) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

var _ pb.OrderServer = &orderServer{}
