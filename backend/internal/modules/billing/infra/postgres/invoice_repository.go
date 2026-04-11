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

func (r *InvoiceRepository) GetByOwnerAndNumber(ctx context.Context, ownerID string, number int) (domain.Invoice, error) {
	var model InvoiceModel
	var status string

	err := r.db.WithContext(ctx).
		Preload("Items").
		Model(&InvoiceModel{}).
		Where("idusuario = ? AND numeracao = ?", ownerID, number).
		First(&model).Error

	if err != nil {
		return domain.Invoice{}, err
	}

	if model.Status == true {
		status = domain.StatusOpen
	} else {
		status = domain.StatusClosed
	}

	items := make([]domain.InvoiceItem, 0, len(model.Items))
	for _, item := range model.Items {
		items = append(items, domain.InvoiceItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	invoice := domain.Invoice{
		OwnerID:   model.OwnerID,
		Number:    model.Numeration,
		Status:    status,
		Items:     items,
		CreatedAt: model.CreatedAt,
	}

	return invoice, nil
}

func (r *InvoiceRepository) Update(ctx context.Context, originalNumber int, invoice domain.Invoice) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model InvoiceModel
		if err := tx.
			Preload("Items").
			Where("idusuario = ? AND numeracao = ?", invoice.OwnerID, originalNumber).
			First(&model).Error; err != nil {
			return err
		}

		model.Numeration = invoice.Number
		model.Status = invoice.Status == domain.StatusOpen

		if err := tx.Model(&model).Updates(map[string]any{
			"numeracao": model.Numeration,
			"status":    model.Status,
		}).Error; err != nil {
			return err
		}

		if err := tx.Where("idnota = ?", model.ID).Delete(&InvoiceItemModel{}).Error; err != nil {
			return err
		}

		items := make([]InvoiceItemModel, 0, len(invoice.Items))
		for _, item := range invoice.Items {
			items = append(items, InvoiceItemModel{
				InvoiceID:   model.ID,
				ProductCode: item.ProductCode,
				Quantity:    item.Quantity,
			})
		}

		if len(items) == 0 {
			return nil
		}

		return tx.Create(&items).Error
	})
}

func (r *InvoiceRepository) UpdateStatus(ctx context.Context, number int, ownerID string, newStatus bool) error {
	err := r.db.WithContext(ctx).
		Model(&InvoiceModel{}).
		Where("idusuario = ? AND numeracao = ?", ownerID, number).
		Update("status", newStatus).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *InvoiceRepository) Delete(ctx context.Context, ownerID string, number int) error {
	return r.db.WithContext(ctx).
		Where("idusuario = ? AND numeracao = ?", ownerID, number).
		Delete(&InvoiceModel{}).Error
}
