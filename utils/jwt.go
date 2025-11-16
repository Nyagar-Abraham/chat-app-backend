package utils

import (
	"os"
	"time"

	"github.com/Nyagar-Abraham/chat-app/models"
	"github.com/golang-jwt/jwt/v4"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
