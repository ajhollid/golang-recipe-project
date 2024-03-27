package models

type JsonIngredient struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

type JsonDirection struct {
	Id        int    `json:"id"`
	Direction string `json:"direction"`
}

type JsonRecipe struct {
	ID          int              `json:"id"`
	Title       string           `json:"title"`
	Ingredients []JsonIngredient `json:"ingredients"`
	Directions  []JsonDirection  `json:"directions"`
}
