package authHandler

import (
	"encoding/json"
	"lockingpostgrescode/auth"
	db "lockingpostgrescode/database"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the structure of the login request body
// @Description The request body for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the structure of the login response
// @Description The response containing the JWT token
type LoginResponse struct {
	Token string `json:"token"`
}

// @Summary User login
// @Description Logs in a user and returns a JWT token
// @Tags AUTH
// @Param loginRequest body authHandler.LoginRequest true "Login Request"
// @Success 200 {object} authHandler.LoginResponse "Login successful, JWT token returned"
// @Failure 400 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials LoginRequest

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, hashedPassword, role, err := validateUser(credentials.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(userID, role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

func validateUser(username string) (int, string, string, error) {
	db, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return 0, "", "", err
	}
	defer db.Close()

	var userID int
	var hashedPassword string
	var roleID int

	err = db.QueryRow("SELECT id, password_hash, role_id FROM users WHERE username = $1", username).Scan(&userID, &hashedPassword, &roleID)
	if err != nil {
		log.Println("User validation failed:", err)
		return 0, "", "", err
	}

	var role string
	err = db.QueryRow("SELECT role_name FROM roles WHERE id = $1", roleID).Scan(&role)
	if err != nil {
		log.Println("Role lookup failed:", err)
		return 0, "", "", err
	}

	return userID, hashedPassword, role, nil
}
