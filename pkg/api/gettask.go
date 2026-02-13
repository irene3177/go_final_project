package api

import (
	"net/http"

	"github.com/irene3177/go_final_project/pkg/db"
)

// getTaskHandler обрабатывает GET запрос для получения задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из query параметра
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		if err.Error() == "Task not found" {
			writeJSONError(w, "Task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "Failed to get task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, task, http.StatusOK)
}
