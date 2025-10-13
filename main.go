package main

import (
	"database/sql"
	"final-project/internal/scheduler"
	"final-project/internal/server"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var Db *sql.DB

func main() {
	logger := log.New(os.Stdout, "SERVER", log.LstdFlags)
	scheduler.DbCheck()
	Db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		logger.Fatal("ошибка создания базы данных", err)
		return
	}
	defer Db.Close()
	err = scheduler.DbCreate(Db)
	if err != nil {
		logger.Fatal("ошибка создания базы данных", err)
		return
	}
	err = os.Setenv("TODO_PASSWORD", "12345")
	if err != nil {
		logger.Fatal("ошибка определения переменной окружения", err)
		return
	}
	server := server.NewServer(logger)
	if err := server.Start(); err != nil {
		logger.Fatal("ошибка при запуске сервера", err.Error())
	}
}
