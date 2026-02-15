package repository

import (
	"database/sql"
	"kasir-api/model"
	"time"
)

type ReportRepository interface {
	GetTodaySummary() (*model.SalesSummary, error)
	GetSummaryByDateRange(startDate, endDate time.Time) (*model.SalesSummary, error)
}

type reportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetTodaySummary() (*model.SalesSummary, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return r.getSummary(startOfDay, startOfDay.Add(24*time.Hour))
}

func (r *reportRepository) GetSummaryByDateRange(startDate, endDate time.Time) (*model.SalesSummary, error) {
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
	return r.getSummary(startDate, endDate)
}

func (r *reportRepository) getSummary(startDate, endDate time.Time) (*model.SalesSummary, error) {
	summary := &model.SalesSummary{}

	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2`,
		startDate, endDate,
	).Scan(&summary.TotalRevenue, &summary.TotalTransactions)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT td.product_id, p.name, SUM(td.quantity) as total_sold
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY td.product_id, p.name
		ORDER BY total_sold DESC
		LIMIT 5`,
		startDate, endDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tp model.TopProduct
		if err := rows.Scan(&tp.ProductID, &tp.ProductName, &tp.TotalSold); err != nil {
			return nil, err
		}
		summary.TopProducts = append(summary.TopProducts, tp)
	}

	if summary.TopProducts == nil {
		summary.TopProducts = []model.TopProduct{}
	}

	return summary, nil
}
