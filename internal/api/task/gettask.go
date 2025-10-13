package task

import (
	"database/sql"
	"errors"
	"final-project/internal/scheduler"
	"net/http"
)

func getTasks(limit int) ([]scheduler.Task, int, error) {
	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer db.Close()
	var tasks []scheduler.Task
	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date ASC LIMIT :limit", sql.Named("limit", limit))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	for rows.Next() {
		var task scheduler.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		tasks = append(tasks, task)
	}
	if tasks == nil {
		return []scheduler.Task{}, http.StatusOK, nil
	}
	return tasks, http.StatusOK, nil
}

func getTask(id string) (scheduler.Task, int, error) {
	var task scheduler.Task
	if id == "" {
		return scheduler.Task{}, http.StatusBadRequest, errors.New("ошибка: требуется указать ID")
	}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return scheduler.Task{}, http.StatusInternalServerError, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	if row.Err() != nil {
		return scheduler.Task{}, http.StatusBadRequest, row.Err()
	}
	row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if len(task.ID) == 0 || task.Title == "" {
		return scheduler.Task{}, http.StatusBadRequest, errors.New("ошибка: задача не найдена")
	}
	return task, http.StatusOK, err
}
