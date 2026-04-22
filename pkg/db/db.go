package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// глобальная переменная для подключения к БД
var DB *sql.DB

// SQL схема для создания таблицы
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT,
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Init инициализирует базу данных
func Init(dbFile string) error {
	// Проверяем, существует ли файл БД
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		if os.IsNotExist(err) {
			install = true
			log.Printf("Database file '%s' not found, will create new database", dbFile)
		} else {
			return fmt.Errorf("Error checking database file: %w", err)
		}
	}

	// Открываем подключение к БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("Failed to open database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("Failed to connect to database: %w", err)
	}

	// Если файла не было, создаем таблицы
	if install {
		log.Printf("Creating database schema...")
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
		log.Printf("Database schema created successfully")
	}

	DB = db
	log.Printf("Database initialized successfully: %s", dbFile)
	return nil
}

// Close закрывает подключение к БД
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
