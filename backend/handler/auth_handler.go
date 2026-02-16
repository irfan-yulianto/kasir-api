package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"kasir-api/model"
	"kasir-api/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		http.Error(w, "Username must be between 3 and 50 characters", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 || len(req.Password) > 100 {
		http.Error(w, "Password must be between 6 and 100 characters", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		req.Role = "kasir"
	}

	resp, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrInvalidRole) {
			http.Error(w, "Role must be 'admin' or 'kasir'", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
