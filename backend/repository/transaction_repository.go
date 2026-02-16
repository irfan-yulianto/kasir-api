package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"kasir-api/model"
)

var ErrInsufficientStock = errors.New("insufficient stock")
var ErrProductNotFound = errors.New("product not found")

type TransactionRepository interface {
	Checkout(items []model.CheckoutItem) (*model.Transaction, error)
	GetAll(startDate, endDate string) ([]model.Transaction, error)
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

		err = tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1 FOR UPDATE", item.ProductID).
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

func (r *transactionRepository) GetAll(startDate, endDate string) ([]model.Transaction, error) {
	query := `SELECT id, total_amount, created_at FROM transactions`
	var args []interface{}

	if startDate != "" && endDate != "" {
		query += ` WHERE created_at >= $1 AND created_at < $2::date + interval '1 day'`
		args = append(args, startDate, endDate)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	var ids []int
	txMap := make(map[int]*model.Transaction)

	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.Details = []model.TransactionDetail{}
		transactions = append(transactions, t)
		ids = append(ids, t.ID)
		txMap[t.ID] = &transactions[len(transactions)-1]
	}

	if len(ids) == 0 {
		return transactions, nil
	}

	// Batch fetch all details in one query instead of N+1
	placeholders := make([]string, len(ids))
	detailArgs := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		detailArgs[i] = id
	}

	detailQuery := fmt.Sprintf(`
		SELECT td.id, td.transaction_id, td.product_id,
			   COALESCE(p.name, 'Deleted Product') as product_name,
			   td.quantity, td.subtotal
		FROM transaction_details td
		LEFT JOIN products p ON td.product_id = p.id
		WHERE td.transaction_id IN (%s)
		ORDER BY td.id`, strings.Join(placeholders, ","))

	detailRows, err := r.db.Query(detailQuery, detailArgs...)
	if err != nil {
		return nil, err
	}
	defer detailRows.Close()

	for detailRows.Next() {
		var d model.TransactionDetail
		if err := detailRows.Scan(&d.ID, &d.TransactionID, &d.ProductID, &d.ProductName, &d.Quantity, &d.Subtotal); err != nil {
			return nil, err
		}
		if tx, ok := txMap[d.TransactionID]; ok {
			tx.Details = append(tx.Details, d)
		}
	}

	return transactions, nil
}
