package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"kasir-api/service"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(service service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleTodayReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var summary interface{}
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, parseErr := time.Parse("2006-01-02", startDateStr)
		if parseErr != nil {
			http.Error(w, "Invalid start_date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		endDate, parseErr := time.Parse("2006-01-02", endDateStr)
		if parseErr != nil {
			http.Error(w, "Invalid end_date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		if startDate.After(endDate) {
			http.Error(w, "start_date must be before end_date", http.StatusBadRequest)
			return
		}

		summary, err = h.service.GetSummaryByDateRange(startDate, endDate)
	} else {
		summary, err = h.service.GetTodaySummary()
	}

	if err != nil {
		http.Error(w, "Failed to get report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
