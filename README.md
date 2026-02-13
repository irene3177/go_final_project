# 📝Планировщик задач (Task Scheduler)

## 📋 Описание проекта

Планировщик задач - это веб-приложение для управления персональными задачами с поддержкой периодических повторений. Проект разработан в рамках финальной работы курса по Golang.


### Основные возможности:

1. ✅ Создание, редактирование и удаление задач
2. 📅 Поддержка различных правил повторения задач:
  * `d N` - повторение каждые N дней
  * `y` - ежегодное повторение
  * `w` - повторение по дням недели
  * `m` - повторение по дням месяца
3. 🔍 Поиск задач по тексту и дате
4. 🔐 Аутентификация по паролю (опционально)
5. 🐳 Запуск в Docker-контейнере

---

## ⭐ Выполненные задания со звёздочкой

### ✅ Шаг 1. Запуск веб-сервера

Поддержка переменной окружения `TODO_PORT` для указания порта

### ✅ Шаг 2. Проектирование БД

Поддержка переменной окружения `TODO_DBFILE` для указания пути к файлу БД

### ✅ Шаг 3. Правила повторения задач

**Полная реализация всех правил:**
  * `d <число>` - повторение через N дней (1-400)
  * `y` - ежегодное повторение
  * `w <1-7,...>` - повторение по дням недели
  * `m <дни> [месяцы]` - повторение по дням месяца
Обработчик `/api/nextdate` для проверки правил

### ✅ Шаг 5. Получение списка задач

* Поиск задач по тексту (в заголовке и комментарии)
* Поиск задач по дате (формат DD.MM.YYYY)

### ✅ Шаг 8. Аутентификация и Docker

* JWT-аутентификация через переменную окружения `TODO_PASSWORD`
* Middleware для защиты API-эндпоинтов
* Docker-образ с поддержкой монтирования volume для БД
* Docker Compose для удобного запуска

---

## 🚀 Запуск локально

### Требования

* Go 1.25 или выше
* SQLite3

### Установка зависимостей

```bash
go mod download
```

### Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `TODO_PORT` | Порт веб-сервера | `7540` |
| `TODO_DBFILE` | Путь к файлу БД | `scheduler.db` |
| `TODO_PASSWORD` | Пароль для аутентификации (опционально) | `""` |
| `TODO_SECRET` | Секретный ключ для JWT | `default-secret-key-change-in-production` |


### Запуск без аутентификации

```bash
go run main.go
```

Приложение будет доступно по адресу: http://localhost:7540

### Запуск с аутентификацией

```bash
export TODO_PASSWORD=mysecretpassword
export TODO_SECRET=my-jwt-secret-key
go run main.go
```

Страница логина: http://localhost:7540/login.html

### Пример .env файла

```env
TODO_PORT=7540
TODO_DBFILE=./data/scheduler.db
TODO_PASSWORD=secure_password_123
TODO_SECRET=your_jwt_secret_key_here
```

---

## 🧪 Запуск тестов

### Настройка tests/settings.go

```go
package tests

var Port = 7540                // Порт для тестов
var DBFile = "../scheduler.db" // Файл БД для тестов
var FullNextDate = true        // true если реализованы все правила (w, m)
var Search = true              // true если реализован поиск
var Token = ""                 // JWT токен (если включена аутентификация)
```

### Запуск отдельных тестов

```bash
# Тест веб-сервера
go test -run ^TestApp$ ./tests

# Тест функции NextDate
go test -run ^TestNextDate$ ./tests

# Тест добавления задачи
go test -run ^TestAddTask$ ./tests

# Тест получения списка задач
go test -run ^TestTasks$ ./tests

# Тест получения задачи по ID
go test -run ^TestTask$ ./tests

# Тест редактирования задачи
go test -run ^TestEditTask$ ./tests

# Тест отметки о выполнении
go test -run ^TestDone$ ./tests

# Тест удаления задачи
go test -run ^TestDelTask$ ./tests
```

### Запуск всех тестов

```bash
go test ./tests
```

---

## 🐳 Запуск в Docker

### Сборка образа

```bash
docker build -t todo-scheduler .
```

### Запуск контейнера

**Без аутентификации:**

```bash
docker run -p 7540:7540 todo-scheduler
```

Или:

```bash
docker run -d \
  --name todo-app \
  -p 7540:7540 \
  -v $(pwd)/data:/data \
  todo-scheduler
```

**С аутентификацией:**

```bash
docker run -d \
  --name todo-app \
  -p 7540:7540 \
  -v $(pwd)/data:/data \
  -e TODO_PASSWORD=mysecretpassword \
  -e TODO_SECRET=my-jwt-secret \
  todo-scheduler
```

**С кастомным портом:**

```bash
docker run -d \
  --name todo-app \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  -e TODO_PORT=8080 \
  todo-scheduler
```

## Запуск через Docker Compose

1. Создайте файл `.env`:

```env
TODO_PASSWORD=mysecretpassword
TODO_SECRET=my-jwt-secret
```

2. Запустите:

```bash
docker-compose up -d
```

3. Остановите:

```bash
docker-compose down
```

---

## 🔌 API Эндпоинты

| Метод | URL | Описание |
|-------|-----|----------|
| `POST` | `/api/task` | Добавить задачу |
| `GET` | `/api/task?id=1` | Получить задачу по ID |
| `PUT` | `/api/task` | Обновить задачу |
| `DELETE` | `/api/task?id=1` | Удалить задачу |
| `GET` | `/api/tasks` | Получить список задач |
| `POST` | `/api/task/done?id=1` | Отметить задачу выполненной |
| `GET` | `/api/nextdate` | Вычислить следующую дату |
| `POST` | `/api/signin` | Аутентификация |


## 🛠 Технологии

* **Backend:** Go 1.25
* **База данных:** SQLite3
* **Аутентификация:** JWT
* **Контейнеризация:** Docker, Docker Compose
* **Фронтенд:** HTML, CSS, JavaScript (готовый)

---

## 📄 Лицензия

Учебный проект, не предназначен для коммерческого использования.

---


