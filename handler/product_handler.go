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

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HandleProducts handles /api/products (GET all, POST)
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProductByID handles /api/products/{id} (GET, PUT, DELETE)
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
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

// HandleProductsByCategory handles /api/categories/{id}/products (GET products by category)
func (h *ProductHandler) HandleProductsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract category ID from path: /api/categories/{id}/products
	path := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	path = strings.TrimSuffix(path, "/products")
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
		products = []model.Product{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) getAll(w http.ResponseWriter, r *http.Request) {
	// Check if include_category query param is set
	includeCategory := r.URL.Query().Get("include_category")

	var products []model.Product
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
		products = []model.Product{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request, id int) {
	// Check if include_category query param is set
	includeCategory := r.URL.Query().Get("include_category")

	var product *model.Product
	var err error

	if includeCategory == "true" {
		product, err = h.service.GetByIDWithCategory(id)
	} else {
		product, err = h.service.GetByID(id)
	}

	if err != nil {
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if product.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&product); err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(id, &product); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted successfully",
	})
}
