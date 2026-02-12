package db

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Role     string    `json:"role"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type File struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	UserId uuid.UUID `json:"user_id"`
}

var Users []*User = []*User{}
var Files []*File = []*File{}
