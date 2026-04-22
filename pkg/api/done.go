package api

import (
	"net/http"
	"time"

	"github.com/irene3177/go_final_project/pkg/db"
)

// doneHandler обрабатывает POST запрос для отметки задачи как выполненной
func doneHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodPost {
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

	// Если нет правила повторения - удаляем задачу
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, "Failed to delete task: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]interface{}{}, http.StatusOK)
		return
	}

	// Есть правило повторения - вычисляем следующую дату
	now := time.Now()
	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, "Failed to calculate next date: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Обновляем дату задачи
	if err := db.UpdateTaskDate(id, nextDate); err != nil {
		writeJSONError(w, "Failed to update task date: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
