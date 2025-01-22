package admin

import (
	"database/sql"
	"net/http"
	"strconv"

	"lockingpostgrescode/exchange"
)

// ExchangeHandler handles currency exchange
// @Summary Exchange currency
// @Tags ADMIN
// @Description Exchange currency from one type to another
// @Param user_id query int true "User ID"
// @Param amount query float64 true "Amount"
// @Security BearerAuth
// @Param from_currency query string true "From Currency (USD/GEL)"
// @Param to_currency query string true "To Currency (USD/GEL)"
// @Success 200 {string} string "Exchange successful"
// @Failure 400 {string} string "Error message"
// @Router /adminexchange [post]
func AdminExchangeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	role := r.Context().Value("role").(string)
	amount, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	fromCurrency := r.URL.Query().Get("from_currency")
	toCurrency := r.URL.Query().Get("to_currency")

	if role == "super_admin" {
		err := exchange.Exchange(db, userID, amount, fromCurrency, toCurrency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
	}

	w.Write([]byte("Exchange successful"))
}
