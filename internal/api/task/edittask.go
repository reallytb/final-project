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

func editTask(w http.ResponseWriter, r *http.Request) (int, error) {
	var task scheduler.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer r.Body.Close()
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return http.StatusBadRequest, err
	}
	if !nextdate.CheckRepeat(task.Repeat) {
		return http.StatusBadRequest, errors.New("ошибка: неверный формат правила повторения")
	}
	if len(task.Title) == 0 {
		return http.StatusBadRequest, errors.New("ошибка: id не может быть пустым")
	}
	if len(task.Title) == 0 {
		return http.StatusBadRequest, errors.New("ошибка: заголовок не может быть пустым")
	}
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		return http.StatusBadRequest, errors.New("ошибка: неверный формат даты")
	}
	defer db.Close()
	res, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return http.StatusBadRequest, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return http.StatusBadRequest, err
	}
	if count == 0 {
		return http.StatusBadRequest, errors.New("ошибка: не изменено ни одного значения")
	}
	return http.StatusOK, nil
}
