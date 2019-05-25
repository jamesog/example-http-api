package database

import "context"

// Storage is the interface for accessing data from a database backend.
type Storage interface {
	Users(ctx context.Context) ([]User, error)
}

// User is an application user.
type User struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"user" db:"user"`
}
