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
// @Param user_id query int true "User ID"
// @Param amount query float64 true "Amount"
// @Security BearerAuth
// @Param fromcurrency query string true "From Currency (USD/GEL)"
// @Param tocurrency query string true "To Currency (USD/GEL)"
// @Success 200 {string} string "Exchange successful"
// @Failure 400 {string} string "Error message"
// @Router /exchange [post]
func ExchangeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	role := r.Context().Value("role").(string)
	amount, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	fromcurrency := r.URL.Query().Get("fromcurrency")
	tocurrency := r.URL.Query().Get("tocurrency")

	if role == "client" {
		_, err := db.Exec(`
			INSERT INTO requests (user_id, operation_type, amount, currency, status, fromcurrency, tocurrency)
			VALUES ($1, 'exchange', $2, 'USD', 'pending',$3, $4)`, userID, amount, fromcurrency, tocurrency)
		if err != nil {
			http.Error(w, "Failed to create exchange request", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("exchange request created successfully"))

	}
}
