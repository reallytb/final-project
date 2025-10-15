package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"final-project/internal/api/nextdate"
	"final-project/internal/scheduler"
)

func addTask(w http.ResponseWriter, r *http.Request) (int64, int, error) {
	var task scheduler.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	defer r.Body.Close()
	if len(task.Title) == 0 {
		err = errors.New("не указан заголовок задачи")
		return 0, http.StatusBadRequest, err
	}
	if task.Date == "" {
		task.Date = time.Now().Format(nextdate.DateFormat)
	}
	if !nextdate.CheckRepeat(task.Repeat) {
		return 0, http.StatusBadRequest, errors.New("ошибка: неверный формат правила повторения")
	}
	t, err := time.Parse(nextdate.DateFormat, task.Date)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	now := time.Now()
	if nextdate.AfterNow(t, now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(nextdate.DateFormat)
		} else {
			next, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, http.StatusBadRequest, err
			}
			task.Date = next
		}

	} else {
		task.Date = now.Format(nextdate.DateFormat)
	}
	id, err := scheduler.AddTask(task)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	return id, http.StatusOK, nil
}
