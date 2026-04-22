package api

import (
	"encoding/json"
	"net/http"

	"github.com/irene3177/go_final_project/pkg/db"
)

// addTaskHandler обрабатывает POST запрос для добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON из запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		writeJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	if err := validateAndFixDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, "Failed to add task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ с ID
	response := map[string]interface{}{
		"id": id,
	}
	writeJSON(w, response, http.StatusCreated)
}
