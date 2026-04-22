package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims структура JWT claims
type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// GenerateToken создает JWT токен
func GenerateToken(password string) (string, error) {
	// Получаем секретный ключ
	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}

	// Создаем claims
	claims := Claims{
		PasswordHash: hashPassword(password),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken проверяет JWT токен
func ValidateToken(tokenString string) (bool, error) {
	// Получаем пароль из переменной окружения
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		// Если пароль не установлен, аутентификация не требуется
		return true, nil
	}

	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return false, fmt.Errorf("Failed to parse token: %w", err)
	}

	// Проверяем валидность токена
	if !token.Valid {
		return false, fmt.Errorf("Invalid token")
	}

	// Получаем claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return false, fmt.Errorf("Invalid claims")
	}

	// Проверяем хэш пароля
	expectedHash := hashPassword(password)
	if claims.PasswordHash != expectedHash {
		return false, fmt.Errorf("Password mismatch")
	}

	return true, nil
}

// AuthMiddleware middleware для проверки аутентификации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, требуется ли аутентификация
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			// Пароль не установлен - пропускаем без проверки
			next(w, r)
			return
		}

		// Получаем токен из куки
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Валидируем токен
		valid, err := ValidateToken(cookie.Value)
		if err != nil || !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// hashPassword создает простой хэш пароля
func hashPassword(password string) string {
	hash := 0
	for _, c := range password {
		hash = (hash*31 + int(c)) % 1000000
	}
	return fmt.Sprintf("%d", hash)
}
