package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"kasir-api/config"
	"kasir-api/database"
	"kasir-api/handler"
	"kasir-api/middleware"
	"kasir-api/repository"
	"kasir-api/service"
)

func main() {
	cfg := config.LoadConfig()

	db := database.NewConnection(cfg.DBConn)
	defer db.Close()

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	reportRepo := repository.NewReportRepository(db)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

	apiKeyMiddleware := middleware.APIKey(cfg.APIKey)

	http.HandleFunc("/api/categories", middleware.CORS(middleware.Logger(categoryHandler.HandleCategories)))
	http.HandleFunc("/api/categories/", middleware.CORS(middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/products") {
			productHandler.HandleProductsByCategory(w, r)
			return
		}
		categoryHandler.HandleCategoryByID(w, r)
	})))

	http.HandleFunc("/api/products", middleware.CORS(middleware.Logger(productHandler.HandleProducts)))
	http.HandleFunc("/api/products/", middleware.CORS(middleware.Logger(apiKeyMiddleware(productHandler.HandleProductByID))))

	http.HandleFunc("/api/checkout", middleware.CORS(middleware.Logger(apiKeyMiddleware(transactionHandler.HandleCheckout))))

	http.HandleFunc("/api/report/today", middleware.CORS(middleware.Logger(reportHandler.HandleTodayReport)))

	http.HandleFunc("/health", middleware.CORS(middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Kasir API Running",
		})
	})))

	addr := "0.0.0.0:" + cfg.Port
	fmt.Printf("Server running on %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
