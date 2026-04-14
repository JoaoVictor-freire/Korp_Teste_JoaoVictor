package http

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/modules/stock/application"
	"korp_backend/internal/modules/stock/domain"
	ai "korp_backend/internal/platform/ai"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/httpx"
)

type Handler struct {
	createProduct application.CreateProductUseCase
	listProducts  application.ListProductsUseCase
	getProduct    application.GetProductUseCase
	updateProduct application.UpdateProductUseCase
	deleteProduct application.DeleteProductUseCase
	decreaseStock application.DecreaseStockUseCase
	aiInsights    application.GenerateAIInsightsUseCase
}

func NewHandler(
	createProduct application.CreateProductUseCase,
	listProducts application.ListProductsUseCase,
	getProduct application.GetProductUseCase,
	updateProduct application.UpdateProductUseCase,
	deleteProduct application.DeleteProductUseCase,
	decreaseStock application.DecreaseStockUseCase,
	aiInsights application.GenerateAIInsightsUseCase,
) Handler {
	return Handler{
		createProduct: createProduct,
		listProducts:  listProducts,
		getProduct:    getProduct,
		updateProduct: updateProduct,
		deleteProduct: deleteProduct,
		decreaseStock: decreaseStock,
		aiInsights:    aiInsights,
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

func (h Handler) GetProduct(c *gin.Context) {
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

	product, err := h.getProduct.Execute(c.Request.Context(), ownerID, code)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrProductCodeRequired):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrProductNotFound):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to get product")
			return
		}
	}

	httpx.JSON(c, http.StatusOK, product)
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

func (h Handler) DecreaseStock(c *gin.Context) {
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

	var request decreaseStockRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.decreaseStock.Execute(c.Request.Context(), application.DecreaseStockInput{
		OwnerID:  ownerID,
		Code:     code,
		Quantity: request.Quantity,
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrProductCodeRequired),
			errors.Is(err, application.ErrDecreaseStockQuantityInvalid):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrProductNotFound):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, domain.ErrInsufficientStock):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to decrease stock")
			return
		}
	}

	httpx.JSON(c, http.StatusOK, product)
}

func (h Handler) GenerateAIInsights(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	log.Printf("ai insights request started: owner_id=%s", ownerID)
	insights, err := h.aiInsights.Execute(c.Request.Context(), ownerID)
	if err != nil {
		log.Printf("ai insights request failed: owner_id=%s err=%v", ownerID, err)
		switch {
		case errors.Is(err, ai.ErrGeminiNotConfigured):
			httpx.Error(c, http.StatusServiceUnavailable, "Servico fora do ar temporariamente")
			return
		default:
			httpx.Error(c, http.StatusBadGateway, "Servico fora do ar temporariamente")
			return
		}
	}

	log.Printf(
		"ai insights request completed: owner_id=%s model=%s products=%d invoices=%d recommendations=%d sources=%d",
		ownerID,
		insights.Model,
		insights.ProductCount,
		insights.InvoiceCount,
		len(insights.BuyRecommendations),
		len(insights.Sources),
	)
	httpx.JSON(c, http.StatusOK, aiInsightsResponse{
		GeneratedAt:        insights.GeneratedAt.Format(time.RFC3339),
		Model:              insights.Model,
		Overview:           insights.Overview,
		Alerts:             insights.Alerts,
		Actions:            insights.Actions,
		BillingNotes:       insights.BillingNotes,
		BuyRecommendations: mapRecommendations(insights.BuyRecommendations),
		SearchQueries:      insights.SearchQueries,
		Sources:            mapSources(insights.Sources),
		ProductCount:       insights.ProductCount,
		InvoiceCount:       insights.InvoiceCount,
		OpenInvoiceCount:   insights.OpenInvoiceCount,
		LowStockCount:      insights.LowStockCount,
		OutOfStockCount:    insights.OutOfStockCount,
	})
}

func mapRecommendations(items []ai.BuyRecommendation) []aiRecommendationDTO {
	result := make([]aiRecommendationDTO, 0, len(items))
	for _, item := range items {
		result = append(result, aiRecommendationDTO{
			Name:          item.Name,
			Category:      item.Category,
			Reason:        item.Reason,
			MarketSignal:  item.MarketSignal,
			StockRelation: item.StockRelation,
		})
	}

	return result
}

func mapSources(items []ai.GroundingSource) []aiSourceDTO {
	result := make([]aiSourceDTO, 0, len(items))
	for _, item := range items {
		result = append(result, aiSourceDTO{
			Title: item.Title,
			URI:   item.URI,
		})
	}

	return result
}
