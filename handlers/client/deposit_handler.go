package client

import (
	"database/sql"
	"net/http"
	"strconv"
)

// DepositHandler handles deposits
// @Summary Deposit funds
// @Tags Client
// @Description Handles deposit requests based on the role of the user
// @Security BearerAuth
// @Param user_id query int true "User ID"
// @Param amount query float64 true "Amount"
// @Param currency query string true "Currency (USD/GEL)"
// @Success 200 {string} string "Deposit successful"
// @Success 201 {string} string "Deposit request created successfully"
// @Failure 400 {string} string "Error message"
// @Failure 403 {string} string "Access denied"
// @Router /deposit [post]
func DepositHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.Context().Value("user_id").(int)
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		http.Error(w, "Currency is required", http.StatusBadRequest)
		return
	}

	if role == "client" {
		_, err := db.Exec(`
			INSERT INTO requests (user_id, operation_type, amount, currency, status)
			VALUES ($1, 'deposit', $2, $3, 'pending')`, userID, amount, currency)
		if err != nil {
			http.Error(w, "Failed to create deposit request", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Deposit request created successfully"))
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
	}
}
