package service

import (
	"kasir-api/model"
	"kasir-api/repository"
	"time"
)

type ReportService interface {
	GetTodaySummary() (*model.SalesSummary, error)
	GetSummaryByDateRange(startDate, endDate time.Time) (*model.SalesSummary, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) GetTodaySummary() (*model.SalesSummary, error) {
	return s.repo.GetTodaySummary()
}

func (s *reportService) GetSummaryByDateRange(startDate, endDate time.Time) (*model.SalesSummary, error) {
	return s.repo.GetSummaryByDateRange(startDate, endDate)
}
