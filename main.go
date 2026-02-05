package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"kasir-api/config"
	"kasir-api/database"
	"kasir-api/handler"
	"kasir-api/repository"
	"kasir-api/service"
)

func main() {
	// ========================================================================
	// LOAD CONFIG
	// ========================================================================
	cfg := config.LoadConfig()

	// ========================================================================
	// CONNECT DATABASE
	// ========================================================================
	db := database.NewConnection(cfg.DBConn)
	defer db.Close()

	// ========================================================================
	// DEPENDENCY INJECTION
	// ========================================================================

	// Category layer
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Product layer
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// Transaction layer
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Report layer
	reportRepo := repository.NewReportRepository(db)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// ========================================================================
	// ROUTES
	// ========================================================================

	// ------------------------------------------------------------------------
	// CATEGORY ROUTES
	// ------------------------------------------------------------------------

	// GET /api/categories - Get all categories
	// POST /api/categories - Create new category
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)

	// GET /api/categories/{id} - Get category detail
	// PUT /api/categories/{id} - Update category
	// DELETE /api/categories/{id} - Delete category
	// Also handles: GET /api/categories/{id}/products - Get products by category
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a request for products by category
		if strings.HasSuffix(r.URL.Path, "/products") {
			productHandler.HandleProductsByCategory(w, r)
			return
		}
		categoryHandler.HandleCategoryByID(w, r)
	})

	// ------------------------------------------------------------------------
	// PRODUCT ROUTES
	// ------------------------------------------------------------------------

	// GET /api/products - Get all products (use ?include_category=true for JOIN)
	// POST /api/products - Create new product
	http.HandleFunc("/api/products", productHandler.HandleProducts)

	// GET /api/products/{id} - Get product detail (use ?include_category=true for JOIN)
	// PUT /api/products/{id} - Update product
	// DELETE /api/products/{id} - Delete product
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// ------------------------------------------------------------------------
	// TRANSACTION ROUTES
	// ------------------------------------------------------------------------

	// POST /api/checkout - Process checkout transaction
	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	// ------------------------------------------------------------------------
	// REPORT ROUTES
	// ------------------------------------------------------------------------

	// GET /api/report/today - Today's sales report
	// Optional query params: ?start_date=2026-02-01&end_date=2026-02-05
	http.HandleFunc("/api/report/today", reportHandler.HandleTodayReport)

	// ------------------------------------------------------------------------
	// HEALTH CHECK
	// ------------------------------------------------------------------------

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Kasir API Running",
		})
	})

	// ========================================================================
	// START SERVER
	// ========================================================================

	addr := "0.0.0.0:" + cfg.Port
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Database connected successfully!")
	fmt.Println("")
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /health")
	fmt.Println("  GET    /api/categories")
	fmt.Println("  POST   /api/categories")
	fmt.Println("  GET    /api/categories/{id}")
	fmt.Println("  PUT    /api/categories/{id}")
	fmt.Println("  DELETE /api/categories/{id}")
	fmt.Println("  GET    /api/categories/{id}/products")
	fmt.Println("  GET    /api/products              (?include_category=true for JOIN)")
	fmt.Println("  POST   /api/products")
	fmt.Println("  GET    /api/products/{id}         (?include_category=true for JOIN)")
	fmt.Println("  PUT    /api/products/{id}")
	fmt.Println("  DELETE /api/products/{id}")
	fmt.Println("  POST   /api/checkout")
	fmt.Println("  GET    /api/report/today")

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
