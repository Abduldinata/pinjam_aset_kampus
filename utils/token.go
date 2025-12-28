package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user_id uint, role string) (string, error) {
	// 1. Ambil secret key dari .env
	secret := os.Getenv("JWT_SECRET")

	// 2. Tentukan isi token (Claims)
	claims := jwt.MapClaims{
		"user_id": user_id,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	// 3. Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Tanda tangani token dengan secret key
	return token.SignedString([]byte(secret))
}