package authHandler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	db "lockingpostgrescode/database"
)

type UserRoleResponse struct {
	Role string `json:"role"`
}

// @Summary Get user role by token
// @Tags AUTH
// @Description Retrieves user role based on provided token
// @Param token header string true "Authorization Token"
// @Success 200 {object} authHandler.UserRoleResponse "User role retrieved successfully"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /role [get]
func RoleHandler(w http.ResponseWriter, r *http.Request) {
	// Extract token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse token
	token := authHeader[7:] // Remove "Bearer " prefix

	// Connect to the database
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var role string
	err = dbConn.QueryRow("SELECT role_name FROM auth_tokens WHERE token = $1", token).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		log.Println("Error retrieving role:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := UserRoleResponse{
		Role: role,
	}

	// Marshal response to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
