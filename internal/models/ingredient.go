package models

type Ingredient struct {
	ID     int
	Recipe Recipe
	Name   string
	Amount string
	Unit   string
}
