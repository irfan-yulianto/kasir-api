package model

// SalesSummary represents daily sales summary report
type SalesSummary struct {
	TotalRevenue      int          `json:"total_revenue"`
	TotalTransactions int          `json:"total_transactions"`
	TopProducts       []TopProduct `json:"top_products"`
}

// TopProduct represents a product with its total sold quantity
type TopProduct struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
}
