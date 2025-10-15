package task

import (
	"errors"
	"net/http"

	"final-project/internal/scheduler"
)

const limit string = "50"

func getTasks() ([]scheduler.Task, int, error) {
	var tasks []scheduler.Task
	tasks, err := scheduler.GetTasks(limit)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return tasks, http.StatusOK, nil
}

func getTask(id string) (scheduler.Task, int, error) {
	if id == "" {
		return scheduler.Task{}, http.StatusBadRequest, errors.New("ошибка: требуется указать ID")
	}
	task, err := scheduler.GetTask(id)
	return task, http.StatusOK, err
}
