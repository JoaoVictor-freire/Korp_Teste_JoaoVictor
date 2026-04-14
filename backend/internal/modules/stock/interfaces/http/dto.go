package http

type createProductRequest struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
}

type decreaseStockRequest struct {
	Quantity int `json:"quantity"`
}

type aiInsightsResponse struct {
	GeneratedAt      string `json:"generated_at"`
	Model            string `json:"model"`
	Content          string `json:"content"`
	ProductCount     int    `json:"product_count"`
	InvoiceCount     int    `json:"invoice_count"`
	OpenInvoiceCount int    `json:"open_invoice_count"`
	LowStockCount    int    `json:"low_stock_count"`
	OutOfStockCount  int    `json:"out_of_stock_count"`
}
