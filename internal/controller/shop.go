package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"github.com/xavimg/Turing/apituringserver/internal/helper"
	"github.com/xavimg/Turing/apituringserver/internal/service"
)

type ShopController interface {
	InsertProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	ListBuys(ctx *gin.Context)
	ListProducts(ctx *gin.Context)
	ProductDetail(ctx *gin.Context)
	AddProductCart(ctx *gin.Context)
	AddCreditCard(ctx *gin.Context)
	DeleteCart(ctx *gin.Context)
	ConfirmPayment(ctx *gin.Context)
}

type shopController struct {
	shopService  service.ShopService
	JWTService   service.JWTService
	adminService service.AdminService
}

func NewShopController(ss service.ShopService, jwt service.JWTService, as service.AdminService) ShopController {
	return &shopController{
		shopService:  ss,
		JWTService:   jwt,
		adminService: as,
	}
}

func (sc *shopController) InsertProduct(ctx *gin.Context) {
	var product entity.Product
	if err := ctx.ShouldBind(&product); err != nil {
		res := helper.BuildErrorResponse(
			"product not add it", err.Error(),
			helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := sc.shopService.InsertProduct(product); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, "product add it")
}

func (sc *shopController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	idP, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "product id is wrong")
		return
	}

	if err := sc.shopService.DeleteProduct(idP); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, "product delete it")
}

func (sc *shopController) UpdateProduct(ctx *gin.Context) {
	var product entity.Product
	if err := ctx.ShouldBind(&product); err != nil {
		res := helper.BuildErrorResponse(
			"product not add it", err.Error(),
			helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	id := ctx.Param("id")
	idP, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "product id is wrong")
		return
	}

	if err := sc.shopService.UpdateProduct(product, idP); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "product update it")

}

func (sc *shopController) ListBuys(ctx *gin.Context) {
	buys, err := sc.shopService.ListBuys()
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, buys)
}

func (sc *shopController) ListProducts(ctx *gin.Context) {
	products, err := sc.shopService.ListProducts()
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (sc *shopController) ProductDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	idP, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "product id is wrong")
		return
	}

	product, err := sc.shopService.ProductDetail(idP)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, product.Detail)
}

func (sc *shopController) AddProductCart(ctx *gin.Context) {
	var product *dto.ProductToCart
	if err := ctx.ShouldBind(&product); err != nil {
		res := helper.BuildErrorResponse(
			"product not add it", err.Error(),
			helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	product, err := sc.shopService.GetProduct(product.Name)
	if err != nil {
		return
	}

	idc := ctx.Param("idcart")
	idCart, err := strconv.Atoi(idc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "product id is wrong")
		return
	}

	if _, err := sc.shopService.StockProduct(product.Name); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	if err := sc.shopService.AddProductCart(idCart, *product); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id := ctx.Param("id")
	idProduct, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "product id is wrong")
		return
	}
	if err := sc.shopService.DeleteProductByName(idProduct, product.Name); err != nil {
		ctx.JSON(http.StatusBadRequest, "error deleting product")
		return
	}

	ctx.JSON(http.StatusOK, "cart created")
}

func (sc *shopController) AddCreditCard(ctx *gin.Context) {
	var product entity.Carts
	if err := ctx.ShouldBind(&product); err != nil {
		res := helper.BuildErrorResponse(
			"product not add it", err.Error(),
			helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	id := ctx.Param("idcart")
	idCart, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "cart id is wrong")
		return
	}

	if _, err := sc.shopService.CartExist(idCart); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	if err := sc.shopService.AddCreditCard(product.CreditCard, idCart); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "payment updated")
}

func (sc *shopController) DeleteCart(ctx *gin.Context) {
	id := ctx.Param("idcart")
	idCart, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "cart id is wrong")
		return
	}

	if _, err := sc.shopService.CartExist(idCart); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	if err := sc.shopService.DeleteCart(idCart); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "cart deleted")
}

func (sc *shopController) ConfirmPayment(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"].(float64)

	idc := ctx.Param("idcart")
	idCart, err := strconv.Atoi(idc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "cart id is wrong")
		return
	}
	if _, err := sc.shopService.CartExist(idCart); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	if err := sc.shopService.ConfirmPayment(id); err != nil {
		ctx.JSON(http.StatusBadRequest, "cart id is wrong")
		return
	}

	ctx.JSON(http.StatusOK, "payment confirmed")
}
