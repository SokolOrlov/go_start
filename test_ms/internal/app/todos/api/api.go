package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	todoService "test_ms/internal/app/todos/service"

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
