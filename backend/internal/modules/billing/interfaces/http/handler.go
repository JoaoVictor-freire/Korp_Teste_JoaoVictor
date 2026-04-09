package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/modules/billing/application"
	"korp_backend/internal/modules/billing/domain"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/httpx"
)

type Handler struct {
	createInvoice application.CreateInvoiceUseCase
	listInvoices  application.ListInvoicesUseCase
}

func NewHandler(
	createInvoice application.CreateInvoiceUseCase,
	listInvoices application.ListInvoicesUseCase,
) Handler {
	return Handler{
		createInvoice: createInvoice,
		listInvoices:  listInvoices,
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
