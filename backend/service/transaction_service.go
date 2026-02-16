package service

import (
	"kasir-api/model"
	"kasir-api/repository"
)

type TransactionService interface {
	Checkout(req model.CheckoutRequest) (*model.Transaction, error)
	GetAll(startDate, endDate string) ([]model.Transaction, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) Checkout(req model.CheckoutRequest) (*model.Transaction, error) {
	return s.repo.Checkout(req.Items)
}

func (s *transactionService) GetAll(startDate, endDate string) ([]model.Transaction, error) {
	return s.repo.GetAll(startDate, endDate)
}
