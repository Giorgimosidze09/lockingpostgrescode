package client

import (
	"database/sql"
	"net/http"
	"strconv"
)

// ExchangeHandler handles currency exchange
// @Summary Exchange currency
// @Tags Client
// @Description Exchange currency from one type to another
// @Security BearerAuth
// @Param amount query float64 true "Amount"
// @Param fromcurrency query string true "From Currency (USD/GEL)"
// @Param tocurrency query string true "To Currency (USD/GEL)"
// @Success 200 {string} string "Exchange successful"
// @Failure 400 {string} string "Error message"
// @Router /exchange [post]
func ExchangeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract user ID and role from context
	userID := r.Context().Value("user_id").(int)
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	fromCurrency := r.URL.Query().Get("fromcurrency")
	if fromCurrency == "" {
		http.Error(w, "From Currency is required", http.StatusBadRequest)
		return
	}

	toCurrency := r.URL.Query().Get("tocurrency")
	if toCurrency == "" {
		http.Error(w, "To Currency is required", http.StatusBadRequest)
		return
	}

	// Only allow clients to create exchange requests
	if role == "client" {
		_, err := db.Exec(`
			INSERT INTO requests (user_id, operation_type, amount, currency, status, fromcurrency, tocurrency)
			VALUES ($1, 'exchange', $2, 'USD', 'pending', $3, $4)`, userID, amount, fromCurrency, toCurrency)
		if err != nil {
			http.Error(w, "Failed to create exchange request", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Exchange request created successfully"))
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
	}
}
