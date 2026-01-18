package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ToDoList represents a task in the to-do list.
type ToDoList struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Task   string             `bson:"task"`
	Status bool               `bson:"status"`
}