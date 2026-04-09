package http

type createProductRequest struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
}
