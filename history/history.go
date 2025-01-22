package history

import (
	"database/sql"
)

// LogHistory logs an operation to the accounts_history table
func LogHistory(tx *sql.Tx, userID int, operationType string) error {
	var usdBalance, usdBlocked, gelBalance, gelBlocked float64

	// Fetch current balances
	err := tx.QueryRow(`
		SELECT usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance 
		FROM accounts 
		WHERE user_id = $1`, userID).Scan(&usdBalance, &usdBlocked, &gelBalance, &gelBlocked)
	if err != nil {
		return err
	}

	// Insert into history table
	_, err = tx.Exec(`
		INSERT INTO accounts_history (user_id, operation_type, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, operationType, usdBalance, usdBlocked, gelBalance, gelBlocked)
	return err
}
