package repository

import (
	"database/sql"
	"errors"
	"kasir-api/model"
)

var ErrInsufficientStock = errors.New("insufficient stock")
var ErrProductNotFound = errors.New("product not found")

type TransactionRepository interface {
	Checkout(items []model.CheckoutItem) (*model.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Checkout(items []model.CheckoutItem) (*model.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var totalAmount int
	var details []model.TransactionDetail

	for _, item := range items {
		var productID int
		var productName string
		var productPrice int
		var productStock int

		err = tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", item.ProductID).
			Scan(&productID, &productName, &productPrice, &productStock)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrProductNotFound
			}
			return nil, err
		}

		if productStock < item.Quantity {
			return nil, ErrInsufficientStock
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).
		Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		var detailID int
		err = tx.QueryRow(
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal,
		).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
		details[i].TransactionID = transactionID
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	var transaction model.Transaction
	err = r.db.QueryRow("SELECT id, total_amount, created_at FROM transactions WHERE id = $1", transactionID).
		Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt)
	if err != nil {
		return nil, err
	}
	transaction.Details = details

	return &transaction, nil
}
