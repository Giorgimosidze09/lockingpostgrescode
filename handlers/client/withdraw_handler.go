package client

import (
	"database/sql"
	"net/http"
	"strconv"
)

// WithdrawHandler handles withdrawals
// @Tags Client
// @Summary Withdraw funds
// @Description Withdraw funds from a user account
// @Security BearerAuth
// @Param amount query float64 true "Amount"
// @Param currency query string true "Currency (USD/GEL)"
// @Success 200 {string} string "Withdrawal successful"
// @Failure 400 {string} string "Error message"
// @Router /withdraw [post]
func WithdrawHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract user ID and role from the token context
	userID := r.Context().Value("user_id").(int)
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Parse amount from query parameters
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// Parse currency from query parameters
	currency := r.URL.Query().Get("currency")
	if currency == "" {
		http.Error(w, "Currency is required", http.StatusBadRequest)
		return
	}

	// Only allow clients to create withdrawal requests
	if role == "client" {
		_, err := db.Exec(`
			INSERT INTO requests (user_id, operation_type, amount, currency, status)
			VALUES ($1, 'withdraw', $2, $3, 'pending')`, userID, amount, currency)
		if err != nil {
			http.Error(w, "Failed to create withdrawal request", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Withdrawal request created successfully"))
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
	}
}
