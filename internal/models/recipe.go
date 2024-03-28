package models

import "time"

type Recipe struct {
	ID          int
	Image       string
	Title       string
	UserId      int
	UpdatedAt   time.Time
	CreatedAt   time.Time
	Ingredients []Ingredient
	Directions  []Direction
	User        User
}
