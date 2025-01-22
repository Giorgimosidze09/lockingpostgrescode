package admin

import (
	"database/sql"
	"net/http"
	"strconv"

	"lockingpostgrescode/withdraw"
)

// WithdrawHandler handles withdrawals
// @Tags ADMIN
// @Summary Withdraw funds
// @Description Withdraw funds from a user account
// @Param user_id query int true "User ID"
// @Param amount query float64 true "Amount"
// @Security BearerAuth
// @Param currency query string true "Currency (USD/GEL)"
// @Success 200 {string} string "Withdrawal successful"
// @Failure 400 {string} string "Error message"
// @Router /adminwithdraw [post]
func AdminWithdrawHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	amount, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	currency := r.URL.Query().Get("currency")
	role := r.Context().Value("role").(string)

	if role == "super_admin" {
		err := withdraw.Withdraw(db, userID, amount, currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
	}

	w.Write([]byte("Withdrawal successful"))
}
