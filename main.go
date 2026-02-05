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

	// Produk layer
	produkRepo := repository.NewProdukRepository(db)
	produkService := service.NewProdukService(produkRepo)
	produkHandler := handler.NewProdukHandler(produkService)

	// Transaction layer (Session 3)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Report layer (Session 3)
	reportRepo := repository.NewReportRepository(db)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handler.NewReportHandler(reportService)

	// ========================================================================
	// ROUTES
	// ========================================================================

	// ------------------------------------------------------------------------
	// CATEGORY ROUTES
	// ------------------------------------------------------------------------

	// GET /api/categories - Ambil semua kategori
	// POST /api/categories - Tambah kategori baru
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)

	// GET /api/categories/{id} - Ambil detail kategori
	// PUT /api/categories/{id} - Update kategori
	// DELETE /api/categories/{id} - Hapus kategori
	// Also handles: GET /api/categories/{id}/produk - Ambil produk berdasarkan kategori
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a request for products by category
		if strings.HasSuffix(r.URL.Path, "/produk") {
			produkHandler.HandleProdukByCategory(w, r)
			return
		}
		categoryHandler.HandleCategoryByID(w, r)
	})

	// ------------------------------------------------------------------------
	// PRODUK ROUTES
	// ------------------------------------------------------------------------

	// GET /api/produk - Ambil semua produk (use ?include_category=true for JOIN)
	// POST /api/produk - Tambah produk baru
	http.HandleFunc("/api/produk", produkHandler.HandleProduk)

	// GET /api/produk/{id} - Ambil detail produk (use ?include_category=true for JOIN)
	// PUT /api/produk/{id} - Update produk
	// DELETE /api/produk/{id} - Hapus produk
	http.HandleFunc("/api/produk/", produkHandler.HandleProdukByID)

	// ------------------------------------------------------------------------
	// TRANSACTION ROUTES (Session 3)
	// ------------------------------------------------------------------------

	// POST /api/checkout - Proses checkout transaksi
	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	// ------------------------------------------------------------------------
	// REPORT ROUTES (Session 3)
	// ------------------------------------------------------------------------

	// GET /api/report/hari-ini - Laporan penjualan hari ini
	// Optional query params: ?start_date=2026-02-01&end_date=2026-02-05
	http.HandleFunc("/api/report/hari-ini", reportHandler.HandleTodayReport)

	// ------------------------------------------------------------------------
	// HEALTH CHECK
	// ------------------------------------------------------------------------

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running with Layered Architecture",
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
	fmt.Println("  GET    /api/categories/{id}/produk  (products by category)")
	fmt.Println("  GET    /api/produk                  (?include_category=true for JOIN)")
	fmt.Println("  POST   /api/produk")
	fmt.Println("  GET    /api/produk/{id}             (?include_category=true for JOIN)")
	fmt.Println("  PUT    /api/produk/{id}")
	fmt.Println("  DELETE /api/produk/{id}")
	fmt.Println("  POST   /api/checkout                (Session 3: Checkout)")
	fmt.Println("  GET    /api/report/hari-ini         (Session 3: Sales Report)")

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
