package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("некорректный формат данных: ожидается три элемента (шаги, вид активности, продолжительность)")
	}

	stepsStr := parts[0]
	activity := parts[1]
	durationStr := parts[2]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, errors.New("не удалось преобразовать количество шагов в число")
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, errors.New("не удалось преобразовать продолжительность в time.Duration")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	totalDistanceMeters := float64(steps) * stepLen
	return totalDistanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distKm := distance(steps, height)
	durationHours := duration.Hours()

	return distKm / durationHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}

	meanSpeedVal := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * meanSpeedVal * durationMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight float64, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}

	meanSpeedVal := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * meanSpeedVal * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	distKm := distance(steps, height)
	meanSpeedVal := meanSpeed(steps, height, duration)

	var calories float64
	var errCalories error

	switch activity {
	case "Ходьба":
		calories, errCalories = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, errCalories = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if errCalories != nil {
		return "", errCalories
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), distKm, meanSpeedVal, calories,
	), nil
}
