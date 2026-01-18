package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madadi-jafar/golang-react-todo/middleware"
)

// GetTodos is a placeholder handler. Replace or remove if not needed.
func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "This endpoint is deprecated or unused"}`))
}

func Router() *mux.Router {
	router := mux.NewRouter()

	// Task routes
	router.HandleFunc("/api/tasks", middleware.GetAllTasks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/tasks", middleware.CreateTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/tasks/{id}", middleware.TaskComplete).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/undoTask/{id}", middleware.UndoTask).Methods("PUT", "OPTIONS") // Changed to PUT for semantic correctness
	router.HandleFunc("/api/deleteTask/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/deleteAllTasks", middleware.DeleteAllTasks).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/deleteCompletedTasks", middleware.DeleteCompletedTasks).Methods("DELETE", "OPTIONS")

	// Optional: keep or remove based on need
	router.HandleFunc("/api/todos", GetTodos).Methods("GET", "OPTIONS")

	return router
}