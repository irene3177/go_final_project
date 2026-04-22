package api

import (
	"encoding/json"
	"net/http"

	"github.com/irene3177/go_final_project/pkg/db"
)

// updateTaskHandler обрабатывает PUT запрос для обновления задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodPut {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON из запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Проверяем наличие ID
	if task.ID == "" {
		writeJSONError(w, "Task ID is required", http.StatusBadRequest)
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

	// Обновляем задачу в БД
	if err := db.UpdateTask(&task); err != nil {
		if err.Error() == "Task not found" {
			writeJSONError(w, "Task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "Failed to update task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
