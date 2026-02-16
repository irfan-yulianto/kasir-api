package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	mux := http.NewServeMux()

	// Auth routes (public)
	mux.HandleFunc("/api/auth/login", authHandler.HandleLogin)
	mux.HandleFunc("/api/auth/register", jwtAuth(adminOnly(authHandler.HandleRegister)))

	// Categories: GET all roles, POST/PUT/DELETE admin only
	mux.HandleFunc("/api/categories", jwtAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			adminOnly(categoryHandler.HandleCategories)(w, r)
			return
		}
		categoryHandler.HandleCategories(w, r)
	}))
	mux.HandleFunc("/api/categories/", jwtAuth(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/products") {
			productHandler.HandleProductsByCategory(w, r)
			return
		}
		if r.Method == http.MethodPut || r.Method == http.MethodDelete {
			adminOnly(categoryHandler.HandleCategoryByID)(w, r)
			return
		}
		categoryHandler.HandleCategoryByID(w, r)
	}))

	// Products: GET all roles, POST/PUT/DELETE admin only
	mux.HandleFunc("/api/products", jwtAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			adminOnly(productHandler.HandleProducts)(w, r)
			return
		}
		productHandler.HandleProducts(w, r)
	}))
	mux.HandleFunc("/api/products/", jwtAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut || r.Method == http.MethodDelete {
			adminOnly(productHandler.HandleProductByID)(w, r)
			return
		}
		productHandler.HandleProductByID(w, r)
	}))

	// Checkout: all authenticated roles
	mux.HandleFunc("/api/checkout", jwtAuth(transactionHandler.HandleCheckout))

	// Transactions: all authenticated roles
	mux.HandleFunc("/api/transactions", jwtAuth(transactionHandler.HandleTransactions))

	// Reports: all authenticated roles
	mux.HandleFunc("/api/report/today", jwtAuth(reportHandler.HandleTodayReport))

	// Health check (public, with DB ping)
	mux.HandleFunc("/health", healthHandler(db))

	// Serve frontend static files
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		fs := http.FileServer(http.Dir(staticDir))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := staticDir + r.URL.Path
			if _, err := os.Stat(path); err != nil {
				http.ServeFile(w, r, staticDir+"/index.html")
				return
			}
			fs.ServeHTTP(w, r)
		})
	}

	// Global middleware stack (outermost → innermost)
	rateLimiter := middleware.NewRateLimiter(10, 30) // 10 req/s, burst 30

	var h http.Handler = mux
	h = middleware.JSONErrors(h)
	h = middleware.MaxBody(1 << 20)(h) // 1MB
	h = middleware.Recovery(h)
	h = rateLimiter.Middleware(h)
	h = middleware.Logger(h)
	h = middleware.RequestID(h)
	h = middleware.CORS(cfg.AllowedOrigins)(h)

	addr := "0.0.0.0:" + cfg.Port
	server := &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("server started", "addr", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	slog.Info("server stopped")
}

func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := "OK"
		dbStatus := "connected"
		httpCode := http.StatusOK

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			status = "degraded"
			dbStatus = "disconnected"
			httpCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		json.NewEncoder(w).Encode(map[string]string{
			"status":   status,
			"database": dbStatus,
		})
	}
}
