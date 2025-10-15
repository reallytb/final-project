package main

import (
	"log"
	"os"

	_ "modernc.org/sqlite"

	"final-project/internal/scheduler"
	"final-project/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "SERVER", log.LstdFlags)
	scheduler.DbCheck()
	err := scheduler.DbCreate()
	if err != nil {
		logger.Fatal("ошибка создания базы данных", err)
		return
	}
	defer scheduler.DB.Close()
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
