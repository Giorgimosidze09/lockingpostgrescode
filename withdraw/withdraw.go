package withdraw

import (
	"database/sql"
	"fmt"
	"time"
)

func Withdraw(db *sql.DB, userID int, amount float64, currency string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	var balance float64
	var updateQuery string

	if currency == "USD" {
		err = tx.QueryRow("SELECT usd_balance FROM accounts WHERE user_id = $1 FOR UPDATE", userID).Scan(&balance)
		updateQuery = "UPDATE accounts SET usd_balance = usd_balance - $1 WHERE user_id = $2"
	} else if currency == "GEL" {
		err = tx.QueryRow("SELECT gel_balance FROM accounts WHERE user_id = $1 FOR UPDATE", userID).Scan(&balance)
		updateQuery = "UPDATE accounts SET gel_balance = gel_balance - $1 WHERE user_id = $2"
	} else {
		return fmt.Errorf("unsupported currency: %s", currency)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch balance: %v", err)
	}

	if balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	_, err = tx.Exec(updateQuery, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %v", err)
	}

	// Log the operation
	_, err = tx.Exec(`
INSERT INTO accounts_history (user_id, operation_type, amount, currency, status, created_at, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance)
SELECT user_id, $1, $2, $3, 'approved', $4, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance
FROM accounts WHERE user_id = $5
`, "withdraw", amount, currency, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to log history: %v", err)
	}

	return tx.Commit()
}
