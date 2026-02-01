package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/model"
	"kasir-api/service"
)

type ProdukHandler struct {
	service service.ProdukService
}

func NewProdukHandler(service service.ProdukService) *ProdukHandler {
	return &ProdukHandler{service: service}
}

// HandleProduk handles /api/produk (GET all, POST)
func (h *ProdukHandler) HandleProduk(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProdukByID handles /api/produk/{id} (GET, PUT, DELETE)
func (h *ProdukHandler) HandleProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProdukByCategory handles /api/categories/{id}/produk (GET products by category)
func (h *ProdukHandler) HandleProdukByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract category ID from path: /api/categories/{id}/produk
	path := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	path = strings.TrimSuffix(path, "/produk")
	categoryID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	products, err := h.service.GetByCategoryID(categoryID)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	if products == nil {
		products = []model.Produk{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProdukHandler) getAll(w http.ResponseWriter, r *http.Request) {
	// Check if include_category query param is set
	includeCategory := r.URL.Query().Get("include_category")

	var products []model.Produk
	var err error

	if includeCategory == "true" {
		products, err = h.service.GetAllWithCategory()
	} else {
		products, err = h.service.GetAll()
	}

	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	if products == nil {
		products = []model.Produk{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProdukHandler) getByID(w http.ResponseWriter, r *http.Request, id int) {
	// Check if include_category query param is set
	includeCategory := r.URL.Query().Get("include_category")

	var produk *model.Produk
	var err error

	if includeCategory == "true" {
		produk, err = h.service.GetByIDWithCategory(id)
	} else {
		produk, err = h.service.GetByID(id)
	}

	if err != nil {
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		return
	}

	if produk == nil {
		http.Error(w, "Produk belum ada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

func (h *ProdukHandler) create(w http.ResponseWriter, r *http.Request) {
	var produk model.Produk
	if err := json.NewDecoder(r.Body).Decode(&produk); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if produk.Nama == "" {
		http.Error(w, "Nama is required", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&produk); err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(produk)
}

func (h *ProdukHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var produk model.Produk
	if err := json.NewDecoder(r.Body).Decode(&produk); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(id, &produk); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Produk belum ada", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

func (h *ProdukHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Produk belum ada", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "sukses delete",
	})
}
