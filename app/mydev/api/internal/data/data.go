package data

import (
	gpb "mydev/api/goods/v1"
)

type DataFactory interface {
	//偷懒做法
	Goods() gpb.GoodsClient
	Users() UserData
	//Categorys() CategoryStore
	//Brands() BrandsStore
	//Banner() BannerStore
	//CategoryBrands() GoodsCategoryBrandStore
}
