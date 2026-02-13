package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/irene3177/go_final_project/pkg/db"
)

const dateFormat = "20060102"

// NextDate вычисляет следующую дату для повторяющейся задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("Repeat rule cannot be empty")
	}

	startDate, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("Invalid start date format: %w", err)
	}

	// Нормализуем даты (убираем время)
	startDate = normalizeDate(startDate)
	nowDate := normalizeDate(now)

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("Invalid repeat format")
	}

	ruleType := parts[0]

	switch ruleType {
	case "d":
		return handleDailyRule(startDate, nowDate, parts)
	case "w":
		return handleWeeklyRule(startDate, nowDate, parts)
	case "m":
		return handleMonthlyRule(startDate, nowDate, parts)
	case "y":
		return handleYearlyRule(startDate, nowDate, parts)
	default:
		return "", fmt.Errorf("Unsupported rule type: %s", ruleType)
	}
}

// handleDailyRule обрабатывает правило d <дни>
func handleDailyRule(startDate, nowDate time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid d rule format: expected 'd <days>'")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("Invalid days value: %w", err)
	}

	if days <= 0 || days > 400 {
		return "", fmt.Errorf("Days must be between 1 and 400")
	}

	date := startDate
	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, nowDate) {
			break
		}
	}
	return formatDate(date), nil
}

func handleWeeklyRule(startDate, nowDate time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid w rule format: expected 'w <weekdays>'")
	}

	weekdaysStr := strings.Split(parts[1], ",")
	weekdaysMap := make(map[int]bool)

	for _, wdStr := range weekdaysStr {
		wd, err := strconv.Atoi(strings.TrimSpace(wdStr))
		if err != nil {
			return "", fmt.Errorf("Invalid weekday value: %w", err)
		}

		if wd < 1 || wd > 7 {
			return "", fmt.Errorf("Weekday must be between 1 and 7")
		}

		weekdaysMap[wd] = true
	}

	date := startDate
	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, nowDate) && isWeekdayInMap(date, weekdaysMap) {
			return formatDate(date), nil
		}

		if date.Sub(startDate).Hours()/24 > 365*3 {
			return "", fmt.Errorf("Cannot find next date within 3 years")
		}
	}
}

func handleMonthlyRule(startDate, nowDate time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("Invalid m rule format: expected 'm <days> [months]'")
	}

	// Парсим дни месяца
	daysStr := strings.Split(parts[1], ",")
	daysMap := make(map[int]bool)
	lastDays := false
	preLastDays := false

	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", fmt.Errorf("invalid day value: %w", err)
		}

		if day == -1 {
			lastDays = true
			continue
		}

		if day == -2 {
			preLastDays = true
			continue
		}

		if day < 1 || day > 31 {
			return "", fmt.Errorf("Day must be between 1 and 31, or -1, -2")
		}

		daysMap[day] = true
	}

	// Парсим месяцы (опционально)
	monthsMap := make(map[int]bool)
	if len(parts) > 2 {
		monthsStr := strings.Split(parts[2], ",")
		for _, monthStr := range monthsStr {
			month, err := strconv.Atoi(strings.TrimSpace(monthStr))
			if err != nil {
				return "", fmt.Errorf("Invalid month value: %w", err)
			}

			if month < 1 || month > 12 {
				return "", fmt.Errorf("Month must be between 1 and 12")
			}
			monthsMap[month] = true
		}
	}

	// Ищем ближайшую подходящую дату
	date := startDate
	maxIterations := 365 * 5

	for range maxIterations {
		date = date.AddDate(0, 0, 1)

		if afterNow(date, nowDate) && isDateValidForMonthlyRule(date, daysMap, monthsMap, lastDays, preLastDays) {
			return formatDate(date), nil
		}
	}

	return "", fmt.Errorf("Cannot find next date within 5 years")

}

func handleYearlyRule(startDate, nowDate time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", fmt.Errorf("Invalid y rule format: expected 'y'")
	}

	date := startDate
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, nowDate) {
			break
		}
	}

	return formatDate(date), nil
}

// isWeekdayInMap проверяет, подходит ли день недели
func isWeekdayInMap(date time.Time, weekdaysMap map[int]bool) bool {
	wd := int(date.Weekday())
	if wd == 0 { // Sunday
		wd = 7
	}
	return weekdaysMap[wd]
}

// isDateValidForMonthlyRule проверяет дату для месячного правила
func isDateValidForMonthlyRule(date time.Time, daysMap map[int]bool, monthsMap map[int]bool, lastDays bool, preLastDays bool) bool {
	// Проверяем месяц если указаны месяцы
	if len(monthsMap) > 0 {
		month := int(date.Month())
		if !monthsMap[month] {
			return false
		}
	}

	day := date.Day()

	// Проверяем специальные дни (-1, -2)
	if lastDays && day == daysInMonth(date.Year(), int(date.Month())) {
		return true
	}

	if preLastDays && day == daysInMonth(date.Year(), int(date.Month()))-1 {
		return true
	}

	return daysMap[day]
}

// daysInMonth возвращает количество дней в месяце
func daysInMonth(year, month int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// afterNow проверяет, что date > now (игнорируя время)
func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()

	dateNormalized := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	nowNormalized := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)

	return dateNormalized.After(nowNormalized)
}

// normalizeDate нормализует дату (убирает время)
func normalizeDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

// formatDate преобразует time.Time в строку формата 20060102
func formatDate(date time.Time) string {
	return date.Format(dateFormat)
}

// validateAndFixDate проверяет и корректирует дату задачи
func validateAndFixDate(task *db.Task) error {
	now := time.Now()
	nowDate := now.Format(dateFormat)

	// Если дата не указана, используем сегодняшнюю
	if task.Date == "" {
		task.Date = nowDate
	}

	// Проверяем формат даты
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("Invalid date format: %w", err)
	}

	// Проверяем правило повторения
	var nextDate string
	if task.Repeat != "" {
		nextDate, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("Invalid repeat rule: %w", err)
		}
	}

	// Проверяем, не устарела ли дата
	if afterNow(now, t) {
		if task.Repeat == "" {
			// Без повторения - ставим сегодня
			task.Date = nowDate
		} else {
			// С повторением - ставим следующую дату
			task.Date = nextDate
		}
	}
	return nil
}

// nextDateHandler обработчик для /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")
	nowParam := r.FormValue("now")

	if dateParam == "" || repeatParam == "" {
		http.Error(w, "Missing required parameters: date and repeat", http.StatusBadRequest)
		return
	}

	// Определяем now
	var now time.Time
	if nowParam != "" {
		parsedNow, err := time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, "Invalid now parameter format", http.StatusBadRequest)
			return
		}
		now = parsedNow
	} else {
		now = time.Now()
	}

	// Вычисляем следующую дату
	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
