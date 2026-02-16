package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	authHandler := handler.NewAuthHandler(authService)

	jwtAuth := middleware.JWTAuth(cfg.JWTSecret)
	adminOnly := middleware.RequireRole("admin")

	// Auth routes (public)
	http.HandleFunc("/api/auth/login", middleware.CORS(middleware.Logger(authHandler.HandleLogin)))
	http.HandleFunc("/api/auth/register", middleware.CORS(middleware.Logger(jwtAuth(adminOnly(authHandler.HandleRegister)))))

	// Categories: GET all roles, POST/PUT/DELETE admin only
	http.HandleFunc("/api/categories", middleware.CORS(middleware.Logger(jwtAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			adminOnly(categoryHandler.HandleCategories)(w, r)
			return
		}
		categoryHandler.HandleCategories(w, r)
	}))))
	http.HandleFunc("/api/categories/", middleware.CORS(middleware.Logger(jwtAuth(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/products") {
			productHandler.HandleProductsByCategory(w, r)
			return
		}
		if r.Method == http.MethodPut || r.Method == http.MethodDelete {
			adminOnly(categoryHandler.HandleCategoryByID)(w, r)
			return
		}
		categoryHandler.HandleCategoryByID(w, r)
	}))))

	// Products: GET all roles, POST/PUT/DELETE admin only
	http.HandleFunc("/api/products", middleware.CORS(middleware.Logger(jwtAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			adminOnly(productHandler.HandleProducts)(w, r)
			return
		}
		productHandler.HandleProducts(w, r)
	}))))
	http.HandleFunc("/api/products/", middleware.CORS(middleware.Logger(jwtAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut || r.Method == http.MethodDelete {
			adminOnly(productHandler.HandleProductByID)(w, r)
			return
		}
		productHandler.HandleProductByID(w, r)
	}))))

	// Checkout: all authenticated roles
	http.HandleFunc("/api/checkout", middleware.CORS(middleware.Logger(jwtAuth(transactionHandler.HandleCheckout))))

	// Transactions: all authenticated roles
	http.HandleFunc("/api/transactions", middleware.CORS(middleware.Logger(jwtAuth(transactionHandler.HandleTransactions))))

	// Reports: all authenticated roles
	http.HandleFunc("/api/report/today", middleware.CORS(middleware.Logger(jwtAuth(reportHandler.HandleTodayReport))))

	// Health check (public)
	http.HandleFunc("/health", middleware.CORS(middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Kasir API Running",
		})
	})))

	// Serve frontend static files
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		fs := http.FileServer(http.Dir(staticDir))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := staticDir + r.URL.Path
			if _, err := os.Stat(path); err != nil {
				// SPA fallback: serve index.html for non-file routes
				http.ServeFile(w, r, staticDir+"/index.html")
				return
			}
			fs.ServeHTTP(w, r)
		})
	}

	addr := "0.0.0.0:" + cfg.Port
	fmt.Printf("Server running on %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
