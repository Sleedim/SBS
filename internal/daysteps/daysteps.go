package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("некорректный формат данных: ожидается два элемента (шаги, продолжительность)")
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, errors.New("не удалось преобразовать количество шагов в число")
	}

	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше нуля")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, errors.New("не удалось преобразовать продолжительность в time.Duration")
	}

	// Явная проверка на нулевую продолжительность (требование теста)
	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше нуля")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
		log.Println("количество шагов должно быть больше нуля")
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)
}
