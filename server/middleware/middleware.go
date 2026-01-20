package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	CreateDBInstance()
}

func CreateDBInstance() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Replace "todo_db" and "tasks" with your actual DB and collection names
	db := client.Database("todo_db")
	collection = db.Collection("tasks")

	log.Println("Connected to MongoDB!")
}

// Helper: Enable CORS and set JSON header
func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// GetAllTasks returns all tasks
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var tasks []bson.M
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

// CreateTask adds a new task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	var task struct {
		Task string `json:"task"`
	}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if task.Task == "" {
		http.Error(w, "Task field is required", http.StatusBadRequest)
		return
	}

	result, err := collection.InsertOne(context.TODO(), bson.M{
		"task":   task.Task,
		"status": false,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"insertedId": result.InsertedID,
	})
}

// TaskComplete marks a task as completed
func TaskComplete(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Task marked as complete"})
}

// UndoTask reverts a task to incomplete
func UndoTask(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": false}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Task undone"})
}

// DeleteTask deletes a single task
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted"})
}

// DeleteAllTasks deletes all tasks
func DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	result, err := collection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"deletedCount": result.DeletedCount})
}

// DeleteCompletedTasks deletes only completed tasks
func DeleteCompletedTasks(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	result, err := collection.DeleteMany(context.TODO(), bson.M{"status": true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"deletedCount": result.DeletedCount})
}