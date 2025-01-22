package exchange

import (
	"database/sql"
	"fmt"
	"time"
)

func Exchange(db *sql.DB, userID int, amount float64, fromCurrency string, toCurrency string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	var fromBalance, toBalance float64
	var updateQuery string

	if fromCurrency == "USD" && toCurrency == "GEL" {
		err = tx.QueryRow("SELECT usd_balance, gel_balance FROM accounts WHERE user_id = $1 FOR UPDATE", userID).
			Scan(&fromBalance, &toBalance)
		updateQuery = "UPDATE accounts SET usd_balance = usd_balance - $1, gel_balance = gel_balance + $2 WHERE user_id = $3"
	} else if fromCurrency == "GEL" && toCurrency == "USD" {
		err = tx.QueryRow("SELECT gel_balance, usd_balance FROM accounts WHERE user_id = $1 FOR UPDATE", userID).
			Scan(&fromBalance, &toBalance)
		updateQuery = "UPDATE accounts SET gel_balance = gel_balance - $1, usd_balance = usd_balance + $2 WHERE user_id = $3"
	} else {
		return fmt.Errorf("unsupported currency conversion: %s to %s", fromCurrency, toCurrency)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch balances: %v", err)
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient funds")
	}

	var conversionRate float64
	if fromCurrency == "USD" && toCurrency == "GEL" {
		conversionRate = 2.0
	} else if fromCurrency == "GEL" && toCurrency == "USD" {
		conversionRate = 0.5
	}

	convertedAmount := amount * conversionRate

	_, err = tx.Exec(updateQuery, amount, convertedAmount, userID)
	if err != nil {
		return fmt.Errorf("failed to update balances: %v", err)
	}

	_, err = tx.Exec(`
INSERT INTO accounts_history (user_id, operation_type, amount, currency, status, created_at, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance)
SELECT user_id, 'withdraw', $1, 'USD', 'approved', $2, usd_balance, usd_blocked_balance, gel_balance, gel_blocked_balance
FROM accounts WHERE user_id = $3
`, amount, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to log history: %v", err)
	}

	return tx.Commit()
}
