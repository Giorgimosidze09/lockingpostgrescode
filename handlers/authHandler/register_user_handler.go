package authHandler

import (
	"encoding/json"
	db "lockingpostgrescode/database"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the structure of the registration request body
// @Description The request body for user registration
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}

// @Summary Register a new user
// @Tags AUTH
// @Description Registers a new user by providing username and password
// @Param registerRequest body authHandler.RegisterRequest true "Register Request"
// @Success 201 {string} string "User registered successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (username, password_hash, role_id) VALUES ($1, $2, $3)", req.Username, string(hashedPassword), req.RoleID)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
