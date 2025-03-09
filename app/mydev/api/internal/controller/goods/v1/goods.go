package goods

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	proto "mydev/api/goods/v1"
	"mydev/app/mydev/api/internal/domain/request"
	"mydev/app/mydev/api/internal/service"
	gin2 "mydev/app/pkg/translator/gin"
	"mydev/pkg/common/core"
	"mydev/pkg/log"
)

type goodsController struct {
	srv   service.ServiceFactory
	trans ut.Translator
}

func NewGoodsController(srv service.ServiceFactory, trans ut.Translator) *goodsController {
	return &goodsController{srv: srv, trans: trans}
}
func (gc *goodsController) List(ctx *gin.Context) {
	log.Info("goods list funciont called ...")

	var r request.GoodsFilter
	if err := ctx.ShouldBindQuery(&r); err != nil {
		gin2.HandleValidatorError(ctx, err, gc.trans)
		return
	}

	gfr := proto.GoodsFilterRequest{
		IsNew:       r.IsNew,
		IsHot:       r.IsHot,
		PriceMax:    r.PriceMax,
		PriceMin:    r.PriceMin,
		TopCategory: r.TopCategory,
		Brand:       r.Brand,
		KeyWords:    r.KeyWords,
		Pages:       r.Pages,
		PagePerNums: r.PagePerNums,
	}
	goodsDTOList, err := gc.srv.Goods().List(ctx, &gfr)
	//TODO: bk
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	reMap := map[string]interface{}{
		"total": goodsDTOList.Total,
	}
	goodsList := make([]interface{}, 0)
	for _, value := range goodsDTOList.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList

	core.WriteResponse(ctx, nil, reMap)

}
