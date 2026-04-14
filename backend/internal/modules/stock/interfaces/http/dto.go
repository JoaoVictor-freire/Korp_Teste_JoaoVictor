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
	GeneratedAt        string                `json:"generated_at"`
	Model              string                `json:"model"`
	Overview           string                `json:"overview"`
	Alerts             []string              `json:"alerts"`
	Actions            []string              `json:"actions"`
	BillingNotes       []string              `json:"billing_notes"`
	BuyRecommendations []aiRecommendationDTO `json:"buy_recommendations"`
	SearchQueries      []string              `json:"search_queries"`
	Sources            []aiSourceDTO         `json:"sources"`
	ProductCount       int                   `json:"product_count"`
	InvoiceCount       int                   `json:"invoice_count"`
	OpenInvoiceCount   int                   `json:"open_invoice_count"`
	LowStockCount      int                   `json:"low_stock_count"`
	OutOfStockCount    int                   `json:"out_of_stock_count"`
}

type aiRecommendationDTO struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	Reason        string `json:"reason"`
	MarketSignal  string `json:"market_signal"`
	StockRelation string `json:"stock_relation"`
}

type aiSourceDTO struct {
	Title string `json:"title"`
	URI   string `json:"uri"`
}
