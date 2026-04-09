package postgres

import "time"

type InvoiceModel struct {
	ID         string             `gorm:"column:idnota;type:uuid;default:uuid_generate_v4();primaryKey"`
	OwnerID    string             `gorm:"column:idusuario;type:uuid;not null;index;uniqueIndex:uidx_nota_usuario_numeracao,priority:1"`
	Status     bool               `gorm:"column:status;not null;default:true"`
	Numeration int                `gorm:"column:numeracao;not null;uniqueIndex:uidx_nota_usuario_numeracao,priority:2"`
	CreatedAt  time.Time          `gorm:"column:criadoem;not null;default:now()"`
	Items      []InvoiceItemModel `gorm:"foreignKey:InvoiceID;references:ID;constraint:OnDelete:CASCADE"`
}

type InvoiceItemModel struct {
	ID          string `gorm:"column:idnotaitem;type:uuid;default:uuid_generate_v4();primaryKey"`
	InvoiceID   string `gorm:"column:idnota;type:uuid;not null;index"`
	ProductCode string `gorm:"column:codigoproduto;type:varchar(500);not null"`
	Quantity    int    `gorm:"column:quantidade;not null"`
}

func (InvoiceModel) TableName() string {
	return "nota"
}

func (InvoiceItemModel) TableName() string {
	return "notaitem"
}

func Models() []any {
	return []any{
		&InvoiceModel{},
		&InvoiceItemModel{},
	}
}
