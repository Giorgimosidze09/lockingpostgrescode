package deposit

import (
	"database/sql"
	"fmt"
	"time"
)

func Deposit(db *sql.DB, userID int, amount float64, currency string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the query to update the account balance
	var query string
	if currency == "USD" {
		query = "UPDATE accounts SET usd_balance = usd_balance + $1 WHERE user_id = $2"
	} else if currency == "GEL" {
		query = "UPDATE accounts SET gel_balance = gel_balance + $1 WHERE user_id = $2"
	} else {
		return fmt.Errorf("unsupported currency: %s", currency)
	}

	// Execute the balance update
	_, err = tx.Exec(query, amount, userID)
	if err != nil {
		// Log failed deposit in accounts_history with "failed" status
		_, _ = tx.Exec(`
			INSERT INTO accounts_history (user_id, operation_type, amount, currency, status, created_at, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance)
			SELECT user_id, $1, $2, $3, 'failed', $4, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance
			FROM accounts WHERE user_id = $5
		`, "deposit", amount, currency, time.Now(), userID)
		return fmt.Errorf("failed to update balance: %v", err)
	}

	// If the deposit is successful, log it as "approved"
	_, err = tx.Exec(`
INSERT INTO accounts_history (user_id, operation_type, amount, currency, status, created_at, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance)
SELECT user_id, $1, $2, $3, 'approved', $4, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance
FROM accounts WHERE user_id = $5
`, "deposit", amount, currency, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to log history: %v", err)
	}

	// Commit the transaction
	return tx.Commit()
}
