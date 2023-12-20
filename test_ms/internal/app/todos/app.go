package todos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"test_ms/internal/app/todos/api"
	"test_ms/internal/app/todos/config"
	"test_ms/internal/app/todos/repo"
	"test_ms/internal/app/todos/repo/todos"
	todoService "test_ms/internal/app/todos/service"

	"github.com/gorilla/mux"
)

type App struct {
	service *todoService.Service
	cfg     *config.Config
	srv     *http.Server
	server  *api.Server
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(ctx context.Context, cfg *config.Config) error {

	a.cfg = cfg

	db := repo.NewDB()

	r := todos.NewRepo(db)

	a.service = todoService.NewService(r)

	a.server = api.NewServer(a.service)

	a.newHttpServer()

	return nil
}

func (a *App) newHttpServer() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/todo/{id}", a.server.Get)
	router.HandleFunc("/api/v1/todos", a.server.GetAll)
	router.HandleFunc("/api/v1/add", a.server.Add).Methods("POST").Headers("Content-Type", "application/json")

	a.srv = &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", a.cfg.HTTP.PORT),
	}
}

func (a *App) Start(ctx context.Context) error {
	//старт серверов grpc, http

	go func() {
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on port :%s: %s\n", a.cfg.HTTP.PORT, err.Error())
		}
	}()

	log.Printf(
		"Serving shop start at port :%s\n",
		a.cfg.HTTP.PORT,
	)

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	<-ctx.Done()

	done := make(chan bool)
	log.Printf("Server is shutting down...")

	// остановка приложения, gracefully shutdown
	go func() {
		if err := a.srv.Shutdown(context.Background()); err != nil {
			log.Fatal("Could not gracefully shutdown the server: ", err.Error())
		}

		log.Printf("Server stopped")
		close(done)
	}()

	<-done
	return nil
}
