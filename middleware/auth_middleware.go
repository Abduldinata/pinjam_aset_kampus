package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware mengecek apakah user login via Cookie (Web) atau Header (Postman)
func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Cek Token dari COOKIE (Untuk Browser)
		cookie, err := c.Cookie("token")
		if err == nil {
			tokenString = cookie
		}

		// 2. Jika Cookie kosong, Cek Token dari HEADER (Untuk Postman/API)
		// Header harus format: "Authorization: Bearer <token_panjang_disini>"
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 3. Jika Token tetap kosong, tolak akses
		if tokenString == "" {
			// Jika request dari browser (HTML), lempar ke login
			// Jika request dari Postman (API), beri JSON error
			if c.GetHeader("Accept") == "application/json" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan, silakan login"})
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
			c.Abort()
			return
		}

		// 4. Validasi JWT Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode enkripsi sesuai (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak valid")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// 5. Cek isi Token & Role
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userRole := claims["role"].(string)
			userID := claims["user_id"].(float64)

			// Simpan data user ke context agar bisa dipakai di controller
			c.Set("user_id", uint(userID))
			c.Set("role", userRole)

			// Cek apakah Role sesuai (misal: halaman admin cuma boleh admin)
			if requiredRole != "" && userRole != requiredRole {
				c.String(http.StatusForbidden, "Akses Ditolak: Halaman ini khusus "+requiredRole)
				c.Abort()
				return
			}

			c.Next() // Lanjut ke controller
		} else {
			// Token kadaluarsa atau palsu
			c.String(http.StatusUnauthorized, "Token tidak valid atau sudah kadaluarsa")
			c.Abort()
		}
	}
}