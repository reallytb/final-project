package scheduler

import (
	"database/sql"
	"os"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

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

func DbCreate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "00000000",
    title TEXT NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat CHAR(128) NOT NULL DEFAULT ""
);`)
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS date_index on scheduler(date);")
	if err != nil {
		return err
	}
	return nil
}
