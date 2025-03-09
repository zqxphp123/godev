package v1

import (
	"context"
	v1 "mydev/app/goods/srv/internal/data/v1"
	"mydev/app/goods/srv/internal/domain/dto"
	metav1 "mydev/pkg/common/meta/v1"
)

type bannerService struct {
	data v1.DataFactory
}

func newBannerService(srv *service) BannerSrv {
	return &bannerService{data: srv.data}
}

func (b *bannerService) Create(ctx context.Context, banner *dto.BannerDTO) error {
	return b.data.Banner().Create(ctx, b.data.Begin(), &banner.BannerDO)
}

func (b *bannerService) Delete(ctx context.Context, id uint64) error {
	return b.data.Banner().Delete(ctx, id)
}

func (b *bannerService) Update(ctx context.Context, banner *dto.BannerDTO) error {
	return b.data.Banner().Update(ctx, b.data.Begin(), &banner.BannerDO)
}

func (b *bannerService) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*dto.BannerListDTO, error) {
	list, err := b.data.Banner().List(ctx, opts, orderby)
	if err != nil {
		return nil, err
	}
	rsp := dto.BannerListDTO{TotalCount: list.TotalCount}
	for _, banner := range list.Items {
		rsp.Items = append(rsp.Items, &dto.BannerDTO{BannerDO: *banner})
	}
	return &rsp, nil
}

type BannerSrv interface {
	//新建轮播图
	Create(ctx context.Context, banner *dto.BannerDTO) error
	//删除轮播图
	Delete(ctx context.Context, id uint64) error
	//修改轮播图信息
	Update(ctx context.Context, banner *dto.BannerDTO) error
	//轮播图列表
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*dto.BannerListDTO, error)
}

var _ BannerSrv = &bannerService{}
