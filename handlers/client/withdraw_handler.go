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
// @Param user_id query int true "User ID"
// @Param amount query float64 true "Amount"
// @Security BearerAuth
// @Param currency query string true "Currency (USD/GEL)"
// @Success 200 {string} string "Withdrawal successful"
// @Failure 400 {string} string "Error message"
// @Router /withdraw [post]
func WithdrawHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	amount, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	currency := r.URL.Query().Get("currency")
	role := r.Context().Value("role").(string)

	if role == "client" {
		_, err := db.Exec(`
			INSERT INTO requests (user_id, operation_type, amount, currency, status)
			VALUES ($1, 'withdraw', $2, $3, 'pending')`, userID, amount, currency)
		if err != nil {
			http.Error(w, "Failed to create exchange request", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Withdrawal  request created"))
	}
}
