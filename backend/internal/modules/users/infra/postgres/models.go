package postgres

import "time"

type UserModel struct {
	ID           string    `gorm:"column:idusuario;type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string    `gorm:"column:nome;type:varchar(100);not null"`
	Email        string    `gorm:"column:email;type:varchar(100);not null"`
	PasswordHash string    `gorm:"column:senha;type:varchar(100);not null"`
	CreatedAt    time.Time `gorm:"column:criadoem;not null;default:now()"`
}

func (UserModel) TableName() string {
	return "usuario"
}

func Models() []any {
	return []any{
		&UserModel{},
	}
}
