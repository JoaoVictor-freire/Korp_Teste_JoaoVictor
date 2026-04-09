package domain

import "time"

type Product struct {
	OwnerID     string    `json:"owner_id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}
