package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// HistoryHandler retrieves account history
// @Tags ADMIN
// @Summary Get account history
// @Description Fetch account history for a user
// @Param user_id path int true "User ID"
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Failure 400 {string} string "Error message"
// @Router /history/{user_id} [get]
func HistoryHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["user_id"])
	role := r.Context().Value("role").(string)

	if role == "super_admin" {
		rows, err := db.Query(`
		SELECT user_id, operation_type, created_at, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance 
		FROM accounts_history 
		WHERE user_id = $1 
		ORDER BY created_at`, userID)
		if err != nil {
			http.Error(w, "Failed to fetch history", http.StatusBadRequest)
			return
		}
		defer rows.Close()

		type AccountHistory struct {
			UserID        int       `json:"user_id"`
			OperationType string    `json:"operation_type"`
			CreatedAt     time.Time `json:"created_at"`
			USDBalance    float64   `json:"usd_balance"`
			USDBlocked    float64   `json:"usd_blocked_balance"`
			GELBalance    float64   `json:"gel_balance"`
			GELBlocked    float64   `json:"gel_blocked_balance"`
		}

		var history []AccountHistory
		for rows.Next() {
			var record AccountHistory
			err := rows.Scan(&record.UserID, &record.OperationType, &record.CreatedAt, &record.USDBalance, &record.USDBlocked, &record.GELBalance, &record.GELBlocked)
			if err != nil {
				http.Error(w, "Failed to parse history", http.StatusInternalServerError)
				return
			}
			history = append(history, record)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating through rows", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(history)

	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
	}

}
