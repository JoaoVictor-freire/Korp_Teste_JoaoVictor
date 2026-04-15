package postgres

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"korp_backend/internal/modules/stock/domain"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product domain.Product) error {
	model := ProductModel{
		OwnerID:     product.OwnerID,
		Description: product.Description,
		Stock:       product.Stock,
		Code:        product.Code,
		CreatedAt:   product.CreatedAt,
	}

	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *ProductRepository) ListByOwner(ctx context.Context, ownerID string) ([]domain.Product, error) {
	var models []ProductModel
	err := r.db.WithContext(ctx).
		Where("idusuario = ?", ownerID).
		Order("codigo asc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	products := make([]domain.Product, 0, len(models))
	for _, model := range models {
		products = append(products, domain.Product{
			OwnerID:     model.OwnerID,
			Code:        model.Code,
			Description: model.Description,
			Stock:       model.Stock,
			CreatedAt:   model.CreatedAt,
		})
	}

	return products, nil
}

func (r *ProductRepository) ExistsByOwnerAndCode(ctx context.Context, ownerID string, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&ProductModel{}).
		Where("idusuario = ? AND codigo = ?", ownerID, code).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ProductRepository) GetByOwnerAndCode(ctx context.Context, ownerID string, code string) (domain.Product, error) {
	var model ProductModel
	err := r.db.WithContext(ctx).
		Where("idusuario = ? AND codigo = ?", ownerID, code).
		First(&model).Error
	if err != nil {
		return domain.Product{}, err
	}

	return domain.Product{
		OwnerID:     model.OwnerID,
		Code:        model.Code,
		Description: model.Description,
		Stock:       model.Stock,
		CreatedAt:   model.CreatedAt,
	}, nil
}

func (r *ProductRepository) Update(ctx context.Context, originalCode string, product domain.Product) error {
	return r.db.WithContext(ctx).
		Model(&ProductModel{}).
		Where("idusuario = ? AND codigo = ?", product.OwnerID, originalCode).
		Updates(map[string]any{
			"codigo":    product.Code,
			"descricao": product.Description,
			"saldo":     product.Stock,
		}).Error
}

func (r *ProductRepository) UpdateStock(ctx context.Context, ownerID string, code string, newStock int) error {
	return r.db.WithContext(ctx).
		Model(&ProductModel{}).
		Where("idusuario = ? AND codigo = ?", ownerID, code).
		Update("saldo", newStock).Error
}

func (r *ProductRepository) DecreaseStock(ctx context.Context, ownerID string, code string, quantity int) (bool, error) {
	result := r.db.WithContext(ctx).
		Model(&ProductModel{}).
		Clauses(clause.Returning{}).
		Where("idusuario = ? AND codigo = ? AND saldo >= ?", ownerID, code, quantity).
		Update("saldo", gorm.Expr("saldo - ?", quantity))

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (r *ProductRepository) Delete(ctx context.Context, ownerID string, code string) error {
	return r.db.WithContext(ctx).
		Where("idusuario = ? AND codigo = ?", ownerID, code).
		Delete(&ProductModel{}).Error
}
