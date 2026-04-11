package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/modules/billing/application"
	"korp_backend/internal/modules/billing/domain"
	stockdomain "korp_backend/internal/modules/stock/domain"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/httpx"
)

type Handler struct {
	createInvoice application.CreateInvoiceUseCase
	listInvoices  application.ListInvoicesUseCase
	closeInvoice  application.CloseInvoiceUseCase
}

func NewHandler(
	createInvoice application.CreateInvoiceUseCase,
	listInvoices application.ListInvoicesUseCase,
	closeInvoice application.CloseInvoiceUseCase,
) Handler {
	return Handler{
		createInvoice: createInvoice,
		listInvoices:  listInvoices,
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

	invoice, err := h.createInvoice.Execute(c.Request.Context(), application.CreateInvoiceInput{
		OwnerID: ownerID,
		Number:  request.Number,
		Items:   items,
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvoiceNumberInvalid),
			errors.Is(err, application.ErrInvoiceItemsRequired),
			errors.Is(err, application.ErrInvoiceItemCodeRequired),
			errors.Is(err, application.ErrInvoiceItemQuantityError):
			httpx.Error(c, http.StatusBadRequest, err.Error())
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

	err = h.closeInvoice.Execute(c.Request.Context(), application.CloseInvoiceInput{
		OwnerID: ownerID,
		Number:  number,
	})
	if err != nil {
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
		case errors.Is(err, domain.ErrInvoiceAlreadyClosed):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		case errors.Is(err, stockdomain.ErrInsufficientStock):
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
