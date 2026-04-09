package postgres

import (
	"context"

	"gorm.io/gorm"

	"korp_backend/internal/modules/billing/domain"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(ctx context.Context, invoice domain.Invoice) error {
	model := InvoiceModel{
		OwnerID:    invoice.OwnerID,
		Status:     invoice.Status == domain.StatusOpen,
		Numeration: invoice.Number,
		CreatedAt:  invoice.CreatedAt,
	}

	model.Items = make([]InvoiceItemModel, 0, len(invoice.Items))
	for _, item := range invoice.Items {
		model.Items = append(model.Items, InvoiceItemModel{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *InvoiceRepository) ListByOwner(ctx context.Context, ownerID string) ([]domain.Invoice, error) {
	var models []InvoiceModel
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("idusuario = ?", ownerID).
		Order("numeracao asc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	invoices := make([]domain.Invoice, 0, len(models))
	for _, model := range models {
		status := domain.StatusClosed
		if model.Status {
			status = domain.StatusOpen
		}

		items := make([]domain.InvoiceItem, 0, len(model.Items))
		for _, item := range model.Items {
			items = append(items, domain.InvoiceItem{
				ProductCode: item.ProductCode,
				Quantity:    item.Quantity,
			})
		}

		invoices = append(invoices, domain.Invoice{
			OwnerID:   model.OwnerID,
			Number:    model.Numeration,
			Status:    status,
			Items:     items,
			CreatedAt: model.CreatedAt,
		})
	}

	return invoices, nil
}

func (r *InvoiceRepository) ExistsByOwnerAndNumber(ctx context.Context, ownerID string, number int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&InvoiceModel{}).
		Where("idusuario = ? AND numeracao = ?", ownerID, number).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
