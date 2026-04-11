package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/modules/stock/application"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/httpx"
)

type Handler struct {
	createProduct application.CreateProductUseCase
	listProducts  application.ListProductsUseCase
	updateProduct application.UpdateProductUseCase
	deleteProduct application.DeleteProductUseCase
}

func NewHandler(
	createProduct application.CreateProductUseCase,
	listProducts application.ListProductsUseCase,
	updateProduct application.UpdateProductUseCase,
	deleteProduct application.DeleteProductUseCase,
) Handler {
	return Handler{
		createProduct: createProduct,
		listProducts:  listProducts,
		updateProduct: updateProduct,
		deleteProduct: deleteProduct,
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

func (h Handler) UpdateProduct(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	originalCode := strings.TrimSpace(c.Param("code"))
	if originalCode == "" {
		httpx.Error(c, http.StatusBadRequest, "product code is required")
		return
	}

	var request createProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.updateProduct.Execute(c.Request.Context(), application.UpdateProductInput{
		OwnerID:      ownerID,
		OriginalCode: originalCode,
		Code:         request.Code,
		Description:  request.Description,
		Stock:        request.Stock,
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrProductCodeRequired),
			errors.Is(err, application.ErrProductDescriptionRequired),
			errors.Is(err, application.ErrProductStockInvalid):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrProductNotFound):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, application.ErrProductAlreadyExists):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to update product")
			return
		}
	}

	httpx.JSON(c, http.StatusOK, product)
}

func (h Handler) DeleteProduct(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		httpx.Error(c, http.StatusBadRequest, "product code is required")
		return
	}

	if err := h.deleteProduct.Execute(c.Request.Context(), ownerID, code); err != nil {
		switch {
		case errors.Is(err, application.ErrProductCodeRequired):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrProductNotFound):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to delete product")
			return
		}
	}

	httpx.JSON(c, http.StatusOK, gin.H{
		"message": "product deleted successfully",
	})
}
