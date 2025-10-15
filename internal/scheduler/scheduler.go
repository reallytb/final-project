package scheduler

import (
	"database/sql"
	"errors"
	"os"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

var DB *sql.DB

func DbCheck() {

	dbFile := "scheduler.db"
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	if install {
		os.Create("scheduler.db")
	}
}

func DbCreate() error {
	var err error
	DB, err = sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "00000000",
    title TEXT NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat CHAR(128) NOT NULL DEFAULT ""
);`)
	if err != nil {
		return err
	}
	_, err = DB.Exec("CREATE INDEX IF NOT EXISTS date_index on scheduler(date);")
	if err != nil {
		return err
	}
	return nil
}

func AddTask(task Task) (int64, error) {
	res, err := DB.Exec("insert into scheduler (date, title, comment, repeat) values (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func EditTask(task Task) error {
	res, err := DB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("ошибка: не изменено ни одного значения")
	}
	return nil
}

func GetTasks(limit string) ([]Task, error) {
	var tasks []Task
	rows, err := DB.Query("SELECT * FROM scheduler ORDER BY date ASC LIMIT :limit", sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if tasks == nil {
		return []Task{}, nil
	}
	return tasks, nil
}

func GetTask(id string) (Task, error) {
	var task Task
	row := DB.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	if row.Err() != nil {
		return Task{}, row.Err()
	}
	row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if len(task.ID) == 0 || task.Title == "" {
		return Task{}, errors.New("ошибка: задача не найдена")
	}
	return task, nil
}
