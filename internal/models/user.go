package models

import "time"

type User struct {
	ID        int
	Image     []byte
	Name      string
	Email     string
	Password  string
	Roles     map[string]Role
	Recipes   []Recipe
	CreatedAt time.Time
	UpdatedAt time.Time
}
