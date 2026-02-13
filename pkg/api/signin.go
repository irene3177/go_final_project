package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/irene3177/go_final_project/pkg/auth"
)

// SigninRequest структура запроса на аутентификацию
type SigninRequest struct {
	Password string `json:"password"`
}

// SigninResponse структура успешного ответа
type SigninResponse struct {
	Token string `json:"token"`
}

// signinHandler обрабатывает POST запрос для аутентификации
func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем запрос
	var req SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Получаем пароль из переменной окружения
	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		writeJSONError(w, "Password not configured", http.StatusInternalServerError)
		return
	}

	// Проверяем пароль
	if req.Password != expectedPassword {
		writeJSONError(w, "Invalid Password", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := auth.GenerateToken(req.Password)
	if err != nil {
		writeJSONError(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   8 * 3600, // 8 hours
	})

	// Возвращаем токен
	response := SigninResponse{
		Token: token,
	}

	writeJSON(w, response, http.StatusOK)
}
