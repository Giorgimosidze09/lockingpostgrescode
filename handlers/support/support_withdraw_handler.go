package support

import (
	"database/sql"
	"lockingpostgrescode/withdraw"
	"log"
	"net/http"
	"strconv"
)

// AcceptRejectDepositRequest handles accepting or rejecting deposit requests
// @Tags Support
// @Summary Support accepts or rejects a deposit request
// @Description Support can either accept or reject a pending deposit request
// @Security BearerAuth
// @Param request_id query int true "Request ID"
// @Param action query string true "Action (accept/reject)"
// @Success 200 {string} string "Request status updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 403 {string} string "Access denied"
// @Router /withdraw/handle [post]
func AcceptRejectWithdrawRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.Context().Value("user_id").(int)
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	if role != "support" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	requestID, err := strconv.Atoi(r.URL.Query().Get("request_id"))
	if err != nil || requestID <= 0 {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	action := r.URL.Query().Get("action")
	if action != "accept" && action != "reject" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	var clientUserID int
	var amount float64
	var currency string

	err = db.QueryRow(`
		SELECT user_id, amount, currency 
		FROM requests 
		WHERE id = $1 AND status = 'pending'`, requestID).Scan(&clientUserID, &amount, &currency)
	if err != nil {
		http.Error(w, "Request not found or already processed", http.StatusNotFound)
		log.Println("Error fetching request details:", err)
		return
	}

	var status string
	if action == "accept" {
		status = "accepted"
	} else {
		status = "rejected"
	}

	_, err = db.Exec(`
		UPDATE requests
		SET status = $1, processed_by = $2
		WHERE id = $3 AND status = 'pending'`, status, userID, requestID)
	if err != nil {
		http.Error(w, "Failed to update request status", http.StatusInternalServerError)
		log.Println("Error updating request status:", err)
		return
	}

	err = withdraw.Withdraw(db, clientUserID, amount, currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request status updated successfully"))
}
