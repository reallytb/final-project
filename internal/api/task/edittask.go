package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"final-project/internal/api/nextdate"
	"final-project/internal/scheduler"
)

func editTask(w http.ResponseWriter, r *http.Request) (int, error) {
	var task scheduler.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer r.Body.Close()
	if !nextdate.CheckRepeat(task.Repeat) {
		return http.StatusBadRequest, errors.New("ошибка: неверный формат правила повторения")
	}
	if len(task.ID) == 0 {
		return http.StatusBadRequest, errors.New("ошибка: id не может быть пустым")
	}
	if len(task.Title) == 0 {
		return http.StatusBadRequest, errors.New("ошибка: заголовок не может быть пустым")
	}
	_, err = time.Parse(nextdate.DateFormat, task.Date)
	if err != nil {
		return http.StatusBadRequest, errors.New("ошибка: неверный формат даты")
	}
	err = scheduler.EditTask(task)
	if err != nil {
		return http.StatusBadRequest, errors.New("ошибка: неверный формат даты")
	}
	return http.StatusOK, nil
}
