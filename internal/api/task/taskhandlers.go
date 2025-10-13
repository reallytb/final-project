package task

import (
	"database/sql"
	"encoding/json"
	"final-project/internal/api/nextdate"
	"final-project/internal/scheduler"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateResponse struct {
	ID string `json:"id"`
}

type taskResp struct {
	Tasks []scheduler.Task `json:"tasks"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		id, statusCode, err := addTask(w, r)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), statusCode)
			return
		}
		var task scheduler.Task
		idString := strconv.Itoa(int(id))
		task.ID = idString
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(CreateResponse{ID: strconv.FormatInt(id, 10)})
	case http.MethodGet:
		r.ParseForm()
		id := r.FormValue("id")
		task, statusCode, err := getTask(id)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(task)
	case http.MethodDelete:
		r.ParseForm()
		id := r.FormValue("id")
		if len(id) == 0 {
			log.Println("ошибка: не указан id")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, "ошибка: не указан id", http.StatusInternalServerError)
			return
		}
		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		res, err := db.Exec("DELETE FROM scheduler where id = :id", sql.Named("id", id))
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if count == 0 {
			log.Println("ошибка: не найдено задачи с заданным id")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, "ошибка: не найдено задачи с заданным id", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	case http.MethodPut:
		statusCode, err := editTask(w, r)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	var task scheduler.Task
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ошибка: задача не найдена")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, "ошибка: задача не найдена", http.StatusBadRequest)
			return
		}
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(task.ID) == 0 {
		log.Println("ошибка: задача не найдена")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		SendJSONError(w, "ошибка: задача не найдена", http.StatusInternalServerError)
		return
	}
	if len(task.Repeat) == 0 {
		_, err := db.Exec("DELETE FROM scheduler where id = :id", sql.Named("id", id))
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	} else {
		nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := db.Exec("UPDATE scheduler SET date = :date where id = :id", sql.Named("date", nextDate), sql.Named("id", id))
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if count == 0 {
			log.Println("ошибка: дата не была изменена")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			SendJSONError(w, "ошибка: дата не была изменена", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, statusCode, err := getTasks(50)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		SendJSONError(w, err.Error(), statusCode)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(taskResp{Tasks: tasks})
}

func SendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
