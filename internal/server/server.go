package server

import (
	"final-project/internal/api/nextdate"
	"final-project/internal/api/signin"
	"final-project/internal/api/task"
	"log"
	"net/http"
	"time"
)

type Server struct {
	logger *log.Logger
	server *http.Server
}

// func NewAddr() (string, int) {
// 	var ServerPort string
// 	fmt.Println("Введите желаемый порт")
// 	fmt.Scan(&ServerPort)
// 	addr := fmt.Sprint(":" + ServerPort)
// 	portInt, _ := strconv.Atoi(ServerPort)
// 	return addr, portInt
// }

func NewServer(logger *log.Logger) *Server {
	addr := ":7540"
	router := newRouter()
	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:         addr,
			Handler:      router,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}
}

func newRouter() http.Handler {
	mux := http.NewServeMux()
	webDir := "./web"
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", nextdate.NexDateHandler)
	mux.HandleFunc("/api/task", signin.Auth(task.TaskHandler))
	mux.HandleFunc("/api/tasks", signin.Auth(task.GetTasksHandler))
	mux.HandleFunc("/api/task/done", signin.Auth(task.TaskDoneHandler))
	mux.HandleFunc("/api/signin", signin.SigninHandler)
	return mux
}

func (s *Server) Start() error {
	s.logger.Printf("Сервер запускается по адресу: %s", s.server.Addr)
	return s.server.ListenAndServe()
}
