package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"kasir-api/model"
	"kasir-api/repository"
	"kasir-api/service"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HandleCheckout handles POST /api/checkout
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if len(req.Items) == 0 {
		http.Error(w, "Items tidak boleh kosong", http.StatusBadRequest)
		return
	}

	for _, item := range req.Items {
		if item.ProductID <= 0 {
			http.Error(w, "product_id harus valid", http.StatusBadRequest)
			return
		}
		if item.Quantity <= 0 {
			http.Error(w, "quantity harus lebih dari 0", http.StatusBadRequest)
			return
		}
	}

	// Process checkout
	transaction, err := h.service.Checkout(req)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientStock) {
			http.Error(w, "Stok tidak mencukupi", http.StatusBadRequest)
			return
		}
		if errors.Is(err, repository.ErrProductNotFound) {
			http.Error(w, "Produk tidak ditemukan", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to process checkout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
