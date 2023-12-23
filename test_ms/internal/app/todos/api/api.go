package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	todoService "test_ms/internal/app/todos/service"
	"test_ms/internal/models"

	"github.com/gorilla/mux"
)

type Server struct {
	todoService *todoService.Service
}

func NewServer(service *todoService.Service) *Server {
	return &Server{
		todoService: service,
	}
}

func (s *Server) GetAll(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	todos, err := s.todoService.GetAll(ctx)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		return
	}
}

func (s *Server) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	vars := mux.Vars(req)
	idStr, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
	}

	id, _ := strconv.Atoi(idStr)

	product, err := s.todoService.Get(ctx, int(id))
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		return
	}
}

func (s *Server) Add(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var t Todo

	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.todoService.Add(ctx, &models.Todo{Task: t.Task, Complete: t.Complete})
	if err != nil {
		return
	}

	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}
