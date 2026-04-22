package api

import (
	"encoding/json"
	"net/http"

	"github.com/irene3177/go_final_project/pkg/auth"
	"github.com/irene3177/go_final_project/pkg/db"
)

func Init() {
	// Обработчик для /api/signin (без аутентификации)
	http.HandleFunc("/api/signin", signinHandler)

	// Обработчик для /api/nextdate (без аутентификации - используется для проверки правил)
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// Обработчики с аутентификацией

	// Обработчик для /api/task (POST, GET, PUT, DELETE)
	http.HandleFunc("/api/task", auth.AuthMiddleware(taskHandler))
	// Обработчик для /api/tasks (GET - список задач)
	http.HandleFunc("/api/tasks", auth.AuthMiddleware(tasksHandler))
	// Обработчик для /api/task/done (POST - отметить выполненной)
	http.HandleFunc("/api/task/done", auth.AuthMiddleware(doneHandler))
}

// taskHandler распределяет запросы по методам
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// deleteTaskHandler обрабатывает DELETE запрос для удаления задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из query параметра
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	// Удаляем задачу из БД
	if err := db.DeleteTask(id); err != nil {
		if err.Error() == "Task not found" {
			writeJSONError(w, "Task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "Failed to delete task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

// writeJSON отправляет JSON ответ
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// writeJSONError отправляет JSON ответ с ошибкой
func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := map[string]string{
		"error": errorMsg,
	}
	writeJSON(w, response, statusCode)
}
