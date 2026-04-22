package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Task структура для хранения задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в БД
func AddTask(task *Task) (int64, error) {
	// Проверяем, что БД инициализирована
	if DB == nil {
		return 0, fmt.Errorf("Database not initialized")
	}

	// SQL запрос для вставки задачи
	query :=
		`INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
		`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("Failed to add task: %w", err)
	}

	// Получаем ID добавленной записи
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to get last insert id: %w", err)
	}

	return id, nil
}

// GetTasks возвращает список задач с сортировкой по дате
func GetTasks(limit int, search string) ([]Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("Database not initialized")
	}

	var query string
	var args []interface{}

	// Проверяем, является ли search датой в формате DD.MM.YYYY
	if search != "" {
		// Пробуем распарсить как дату
		if t, err := time.Parse("02.01.2006", search); err == nil {
			// Это дата - ищем по точному совпадению
			date := t.Format("20060102")
			query = `
				SELECT id, date, title, comment, repeat
				FROM scheduler
				WHERE date = ?
				ORDER BY date ASC
				LIMIT ?
			`

			args = []interface{}{date, limit}
		} else {
			// Это поиск по тексту
			searchTerm := "%" + search + "%"
			query = `
				SELECT id, date, title, comment, repeat
				FROM scheduler
				WHERE LOWER(title) LIKE ? OR LOWER(comment) LIKE ?
				ORDER BY date ASC
				LIMIT ?
			`

			args = []interface{}{searchTerm, searchTerm, limit}
		}
	} else {
		// Без поиска - просто все задачи
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date ASC
			LIMIT ?
		`

		args = []interface{}{limit}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	// Возвращаем пустой слайс вместо nil для корректного JSON
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, nil
}

// GetTask возвращает задачу по ID
func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("Database not initialized")
	}

	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`

	var task Task
	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Task not found")
		}
		return nil, fmt.Errorf("Failed to get task: %w", err)
	}
	return &task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("Database not initialized")
	}

	// Проверяем, что ID не пустой
	if task.ID == "" {
		return fmt.Errorf("Task ID is required")
	}

	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("Failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("Task not found")
	}
	return nil
}

// UpdateTaskDate обновляет только дату задачи
func UpdateTaskDate(id string, newDate string) error {
	if DB == nil {
		return fmt.Errorf("Database not initialized")
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	result, err := DB.Exec(query, newDate, id)
	if err != nil {
		return fmt.Errorf("Failed to update task date: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("Task not found")
	}

	return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
	if DB == nil {
		return fmt.Errorf("Database not initialized")
	}

	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("Task not found")
	}

	return nil
}

// TaskExists проверяет существование задачи по ID
func TaskExists(id string) (bool, error) {
	if DB == nil {
		return false, fmt.Errorf("Database not initialized")
	}

	query := `SELECT COUNT(*) FROM scheduler WHERE id = ?`

	var count int
	err := DB.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Failed to check existence: %w", err)
	}
	return count > 0, nil
}
