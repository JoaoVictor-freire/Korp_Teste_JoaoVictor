package postgres

import "time"

type ProductModel struct {
	ID          string    `gorm:"column:idproduto;type:uuid;default:uuid_generate_v4();primaryKey"`
	OwnerID     string    `gorm:"column:idusuario;type:uuid;not null;index;uniqueIndex:uidx_produto_usuario_codigo,priority:1"`
	Description string    `gorm:"column:descricao;type:varchar(100);not null"`
	Stock       int       `gorm:"column:saldo;not null"`
	Code        string    `gorm:"column:codigo;type:varchar(500);not null;uniqueIndex:uidx_produto_usuario_codigo,priority:2"`
	CreatedAt   time.Time `gorm:"column:criadoem;not null;default:now()"`
}

func (ProductModel) TableName() string {
	return "produto"
}

func Models() []any {
	return []any{
		&ProductModel{},
	}
}
