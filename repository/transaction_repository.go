package repository

import (
	"database/sql"
	"errors"
	"kasir-api/model"
)

var ErrInsufficientStock = errors.New("stok tidak mencukupi")
var ErrProductNotFound = errors.New("produk tidak ditemukan")

type TransactionRepository interface {
	Checkout(items []model.CheckoutItem) (*model.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Checkout processes a checkout request using database transaction
func (r *transactionRepository) Checkout(items []model.CheckoutItem) (*model.Transaction, error) {
	// Start database transaction
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

	// Process each item
	for _, item := range items {
		// Get product details
		var produkID int
		var produkNama string
		var produkHarga int
		var produkStok int

		query := "SELECT id, nama, harga, stok FROM produk WHERE id = $1"
		err = tx.QueryRow(query, item.ProductID).Scan(&produkID, &produkNama, &produkHarga, &produkStok)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrProductNotFound
			}
			return nil, err
		}

		// Validate stock
		if produkStok < item.Quantity {
			return nil, ErrInsufficientStock
		}

		// Calculate subtotal
		subtotal := produkHarga * item.Quantity
		totalAmount += subtotal

		// Update product stock
		updateQuery := "UPDATE produk SET stok = stok - $1 WHERE id = $2"
		_, err = tx.Exec(updateQuery, item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// Store detail for later insertion
		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: produkNama,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Insert transaction
	var transactionID int
	insertTxQuery := "INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id"
	err = tx.QueryRow(insertTxQuery, totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Insert transaction details
	insertDetailQuery := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id"
	for i := range details {
		var detailID int
		err = tx.QueryRow(insertDetailQuery, transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
		details[i].TransactionID = transactionID
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Get created_at
	var transaction model.Transaction
	getQuery := "SELECT id, total_amount, created_at FROM transactions WHERE id = $1"
	err = r.db.QueryRow(getQuery, transactionID).Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt)
	if err != nil {
		return nil, err
	}
	transaction.Details = details

	return &transaction, nil
}
