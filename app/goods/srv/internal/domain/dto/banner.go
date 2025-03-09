package dto

import "mydev/app/goods/srv/internal/domain/do"

type BannerDTO struct {
	do.BannerDO
}

type BannerListDTO struct {
	TotalCount int64        `json:"totalCount,omitempty"`
	Items      []*BannerDTO `json:"items"`
}
