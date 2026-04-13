package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/modules/billing/application"
	"korp_backend/internal/modules/billing/domain"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/httpx"
)

type Handler struct {
	createInvoice application.CreateInvoiceUseCase
	listInvoices  application.ListInvoicesUseCase
	getInvoice    application.GetInvoiceUseCase
	updateInvoice application.UpdateInvoiceUseCase
	deleteInvoice application.DeleteInvoiceUseCase
	closeInvoice  application.CloseInvoiceUseCase
}

func NewHandler(
	createInvoice application.CreateInvoiceUseCase,
	listInvoices application.ListInvoicesUseCase,
	getInvoice application.GetInvoiceUseCase,
	updateInvoice application.UpdateInvoiceUseCase,
	deleteInvoice application.DeleteInvoiceUseCase,
	closeInvoice application.CloseInvoiceUseCase,
) Handler {
	return Handler{
		createInvoice: createInvoice,
		listInvoices:  listInvoices,
		getInvoice:    getInvoice,
		updateInvoice: updateInvoice,
		deleteInvoice: deleteInvoice,
		closeInvoice:  closeInvoice,
	}
}

func (h Handler) CreateInvoice(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	var request createInvoiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	items := make([]domain.InvoiceItem, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, domain.InvoiceItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	invoice, err := h.createInvoice.Execute(application.WithAuthorizationHeader(c.Request.Context(), c.GetHeader("Authorization")), application.CreateInvoiceInput{
		OwnerID: ownerID,
		Number:  request.Number,
		Items:   items,
	})
	if err != nil {
		var productNotFoundErr application.InvoiceProductNotFoundError
		var outOfStockErr application.InvoiceOutOfStockError
		switch {
		case errors.Is(err, application.ErrInvoiceNumberInvalid),
			errors.Is(err, application.ErrInvoiceItemsRequired),
			errors.Is(err, application.ErrInvoiceItemCodeRequired),
			errors.Is(err, application.ErrInvoiceItemQuantityError):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrStockUnauthorized),
			errors.Is(err, application.ErrStockUnavailable),
			errors.Is(err, application.ErrStockCircuitOpen):
			httpx.Error(c, http.StatusBadGateway, err.Error())
			return
		case errors.As(err, &productNotFoundErr):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		case errors.As(err, &outOfStockErr):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		case errors.Is(err, application.ErrInvoiceAlreadyExists):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to create invoice")
			return
		}
	}

	httpx.JSON(c, http.StatusCreated, invoice)
}

func (h Handler) ListInvoices(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	invoices, err := h.listInvoices.Execute(c.Request.Context(), ownerID)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "failed to list invoices")
		return
	}

	httpx.JSON(c, http.StatusOK, invoices)
}

func (h Handler) GetInvoice(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid invoice number")
		return
	}

	invoice, err := h.getInvoice.Execute(c.Request.Context(), ownerID, number)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvoiceNumberInvalid):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrCloseInvoiceNotFound):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to get invoice")
			return
		}
	}

	httpx.JSON(c, http.StatusOK, invoice)
}

func (h Handler) UpdateInvoice(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid invoice number")
		return
	}

	var request createInvoiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	items := make([]domain.InvoiceItem, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, domain.InvoiceItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	invoice, err := h.updateInvoice.Execute(c.Request.Context(), application.UpdateInvoiceInput{
		OwnerID:        ownerID,
		OriginalNumber: number,
		Number:         request.Number,
		Items:          items,
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvoiceNumberInvalid),
			errors.Is(err, application.ErrInvoiceItemsRequired),
			errors.Is(err, application.ErrInvoiceItemCodeRequired),
			errors.Is(err, application.ErrInvoiceItemQuantityError):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrCloseInvoiceNotFound):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, application.ErrInvoiceAlreadyExists),
			errors.Is(err, application.ErrInvoiceCannotEditClosed):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to update invoice")
			return
		}
	}

	httpx.JSON(c, http.StatusOK, invoice)
}

func (h Handler) DeleteInvoice(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid invoice number")
		return
	}

	if err := h.deleteInvoice.Execute(c.Request.Context(), ownerID, number); err != nil {
		switch {
		case errors.Is(err, application.ErrInvoiceNumberInvalid):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrCloseInvoiceNotFound):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, application.ErrInvoiceCannotDeleteClosed):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to delete invoice")
			return
		}
	}

	httpx.JSON(c, http.StatusOK, gin.H{
		"message": "invoice deleted successfully",
	})
}

func (h Handler) CloseInvoice(c *gin.Context) {
	ownerID, ok := auth.UserIDFromContext(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "missing auth context")
		return
	}

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid invoice number")
		return
	}

	err = h.closeInvoice.Execute(application.WithAuthorizationHeader(c.Request.Context(), c.GetHeader("Authorization")), application.CloseInvoiceInput{
		OwnerID: ownerID,
		Number:  number,
	})
	if err != nil {
		var insufficientStockErr application.InsufficientStockError
		switch {
		case errors.Is(err, application.ErrCloseInvoiceOwnerRequired),
			errors.Is(err, application.ErrInvoiceNumberInvalid):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrCloseInvoiceNotFound):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, application.ErrCloseInvoiceProductNotFound):
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, application.ErrStockUnauthorized),
			errors.Is(err, application.ErrStockUnavailable),
			errors.Is(err, application.ErrStockCircuitOpen):
			httpx.Error(c, http.StatusBadGateway, err.Error())
			return
		case errors.Is(err, domain.ErrInvoiceAlreadyClosed):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		case errors.As(err, &insufficientStockErr):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to close invoice")
			return
		}
	}

	httpx.JSON(c, http.StatusOK, gin.H{
		"message": "invoice closed successfully",
	})
}
