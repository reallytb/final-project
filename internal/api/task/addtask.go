package task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"final-project/internal/api/nextdate"
	"final-project/internal/scheduler"
	"net/http"
	"time"
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
		task.Date = time.Now().Format("20060102")
	}
	if !nextdate.CheckRepeat(task.Repeat) {
		return 0, http.StatusBadRequest, errors.New("ошибка: неверный формат правила повторения")
	}
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	now := time.Now()
	if nextdate.AfterNow(t, now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format("20060102")
		} else {
			next, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, http.StatusBadRequest, err
			}
			task.Date = next
		}

	} else {
		task.Date = now.Format("20060102")
	}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	defer db.Close()
	res, err := db.Exec("insert into scheduler (date, title, comment, repeat) values (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	return id, http.StatusOK, nil
}
