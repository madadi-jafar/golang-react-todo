package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	// Add your routes here
	r.HandleFunc("/api/todos", GetTodos).Methods("GET")
	// ... other handlers
	return r
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("To-do list placeholder"))
}