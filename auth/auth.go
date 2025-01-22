package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims struct to store the JWT claims
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token with userID and role
func GenerateToken(userID int, role string) (string, error) {
	// Set the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "yourapp",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the JWT token and returns the claims if valid
func ValidateToken(tokenString string) (*Claims, error) {
	// Define your secret key (should match the key used to sign the token)
	secretKey := []byte("your_secret_key")

	// Parse the token with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is signed with the correct method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	// Optionally, check if the token is expired
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
