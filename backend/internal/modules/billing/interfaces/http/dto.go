package http

type createInvoiceRequest struct {
	Number int                 `json:"number"`
	Items  []createInvoiceItem `json:"items"`
}

type createInvoiceItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}
