// Code generated by protoc-gen-gin. DO NOT EDIT.

package proto

import (
	gin "github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
)

type GoodsHttpServer struct {
	server GoodsServer
	router gin.IRouter
}

func RegisterGoodsServerHTTPServer(srv GoodsServer, r gin.IRouter) {
	s := GoodsHttpServer{
		server: srv,
		router: r,
	}
	s.RegisterService()
}

func (s *GoodsHttpServer) GoodsList_0(c *gin.Context) {
	var in GoodsFilterRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.GoodsList(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) BatchGetGoods_0(c *gin.Context) {
	var in BatchGoodsIdInfo

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.BatchGetGoods(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) CreateGoods_0(c *gin.Context) {
	var in CreateGoodsInfo

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.CreateGoods(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) DeleteGoods_0(c *gin.Context) {
	var in DeleteGoodsInfo

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.DeleteGoods(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) UpdateGoods_0(c *gin.Context) {
	var in CreateGoodsInfo

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.UpdateGoods(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) GetGoodsDetail_0(c *gin.Context) {
	var in GoodInfoRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.GetGoodsDetail(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) GetAllCategorysList_0(c *gin.Context) {
	var in emptypb.Empty

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.GetAllCategorysList(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) GetSubCategory_0(c *gin.Context) {
	var in CategoryListRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.GetSubCategory(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) CreateCategory_0(c *gin.Context) {
	var in CategoryInfoRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.CreateCategory(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) DeleteCategory_0(c *gin.Context) {
	var in DeleteCategoryRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.DeleteCategory(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) UpdateCategory_0(c *gin.Context) {
	var in CategoryInfoRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.UpdateCategory(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) BrandList_0(c *gin.Context) {
	var in BrandFilterRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.BrandList(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) CreateBrand_0(c *gin.Context) {
	var in BrandRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.CreateBrand(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) DeleteBrand_0(c *gin.Context) {
	var in BrandRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.DeleteBrand(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) UpdateBrand_0(c *gin.Context) {
	var in BrandRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.UpdateBrand(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) BannerList_0(c *gin.Context) {
	var in BannerListReq

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.BannerList(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) CreateBanner_0(c *gin.Context) {
	var in BannerRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.CreateBanner(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) DeleteBanner_0(c *gin.Context) {
	var in BannerRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.DeleteBanner(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) UpdateBanner_0(c *gin.Context) {
	var in BannerRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.UpdateBanner(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) CategoryBrandList_0(c *gin.Context) {
	var in CategoryBrandFilterRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.CategoryBrandList(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) GetCategoryBrandList_0(c *gin.Context) {
	var in CategoryInfoRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.GetCategoryBrandList(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) CreateCategoryBrand_0(c *gin.Context) {
	var in CategoryBrandRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.CreateCategoryBrand(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) DeleteCategoryBrand_0(c *gin.Context) {
	var in CategoryBrandRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.DeleteCategoryBrand(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) UpdateCategoryBrand_0(c *gin.Context) {
	var in CategoryBrandRequest

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := s.server.UpdateCategoryBrand(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (s *GoodsHttpServer) RegisterService() {

	s.router.Handle("POST", "", s.GoodsList_0)

	s.router.Handle("POST", "", s.BatchGetGoods_0)

	s.router.Handle("POST", "", s.CreateGoods_0)

	s.router.Handle("POST", "", s.DeleteGoods_0)

	s.router.Handle("POST", "", s.UpdateGoods_0)

	s.router.Handle("POST", "", s.GetGoodsDetail_0)

	s.router.Handle("POST", "", s.GetAllCategorysList_0)

	s.router.Handle("POST", "", s.GetSubCategory_0)

	s.router.Handle("POST", "", s.CreateCategory_0)

	s.router.Handle("POST", "", s.DeleteCategory_0)

	s.router.Handle("POST", "", s.UpdateCategory_0)

	s.router.Handle("POST", "", s.BrandList_0)

	s.router.Handle("POST", "", s.CreateBrand_0)

	s.router.Handle("POST", "", s.DeleteBrand_0)

	s.router.Handle("POST", "", s.UpdateBrand_0)

	s.router.Handle("POST", "", s.BannerList_0)

	s.router.Handle("POST", "", s.CreateBanner_0)

	s.router.Handle("POST", "", s.DeleteBanner_0)

	s.router.Handle("POST", "", s.UpdateBanner_0)

	s.router.Handle("POST", "", s.CategoryBrandList_0)

	s.router.Handle("POST", "", s.GetCategoryBrandList_0)

	s.router.Handle("POST", "", s.CreateCategoryBrand_0)

	s.router.Handle("POST", "", s.DeleteCategoryBrand_0)

	s.router.Handle("POST", "", s.UpdateCategoryBrand_0)

}
