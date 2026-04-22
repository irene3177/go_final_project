package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/irene3177/go_final_project/pkg/api"
	"github.com/irene3177/go_final_project/pkg/db"
)

func Run() error {
	port := 7540
	webDir := "./web"

	// База данных
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "./data/scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		return err
	}
	defer db.Close()

	// Проверяем настройки аутентификации
	if password := os.Getenv("TODO_PASSWORD"); password != "" {
		log.Println("Authentication enabled")
	} else {
		log.Println("Authentication disabled")
	}

	// Порт
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	// Регистрируем API
	api.Init()

	// Статические файлы
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Starting the server on port %d...", port)
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)

}
