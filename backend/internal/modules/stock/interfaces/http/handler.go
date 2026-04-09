package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/modules/stock/application"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/httpx"
)

type Handler struct {
	createProduct application.CreateProductUseCase
	listProducts  application.ListProductsUseCase
}

func NewHandler(
	createProduct application.CreateProductUseCase,
	listProducts application.ListProductsUseCase,
) Handler {
	return Handler{
		createProduct: createProduct,
		listProducts:  listProducts,
	}
}

func (h Handler) CreateProduct(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	var request createProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.createProduct.Execute(c.Request.Context(), application.CreateProductInput{
		OwnerID:     ownerID,
		Code:        request.Code,
		Description: request.Description,
		Stock:       request.Stock,
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrProductCodeRequired),
			errors.Is(err, application.ErrProductDescriptionRequired),
			errors.Is(err, application.ErrProductStockInvalid):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrProductAlreadyExists):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to create product")
			return
		}
	}

	httpx.JSON(c, http.StatusCreated, product)
}

func (h Handler) ListProducts(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	products, err := h.listProducts.Execute(c.Request.Context(), ownerID)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "failed to list products")
		return
	}

	httpx.JSON(c, http.StatusOK, products)
}
