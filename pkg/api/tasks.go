package api

import (
	"net/http"

	"github.com/irene3177/go_final_project/pkg/db"
)

// TasksResponse структура ответа со списком задач
type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET запрос для получения списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр поиска
	search := r.URL.Query().Get("search")

	// Лимит задач по умолчанию - 50
	limit := 50

	// Получаем задачи из БД
	tasks, err := db.GetTasks(limit, search)
	if err != nil {
		writeJSONError(w, "Failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	response := TasksResponse{
		Tasks: tasks,
	}
	writeJSON(w, response, http.StatusOK)
}
