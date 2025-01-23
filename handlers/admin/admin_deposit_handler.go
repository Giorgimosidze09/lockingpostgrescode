package admin

import (
	"database/sql"
	"lockingpostgrescode/deposit"
	"net/http"
	"strconv"
)

// AdminDepositHandler handles deposits by admin or super_admin
// @Summary Deposit funds
// @Tags ADMIN
// @Description Handles deposit requests based on the role of the user
// @Security BearerAuth
// @Param user_id query int true "Client User ID"
// @Param amount query float64 true "Amount"
// @Param currency query string true "Currency (USD/GEL)"
// @Success 200 {string} string "Deposit successful"
// @Success 201 {string} string "Deposit request created successfully"
// @Failure 400 {string} string "Error message"
// @Failure 403 {string} string "Access denied"
// @Router /admindeposit [post]
func AdminDepositHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || userID == 0 {
		http.Error(w, "Invalid client user ID", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		http.Error(w, "Currency is required", http.StatusBadRequest)
		return
	}

	if role == "super_admin" {
		err = deposit.Deposit(db, userID, amount, currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("Deposit successful"))
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
	}
}
