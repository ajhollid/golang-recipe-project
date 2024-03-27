package repository

import (
	"github.com/popnfresh234/recipe-app-golang/internal/models"
)

type DatabaseRepo interface {
	GetAllRecipes() ([]models.Recipe, error)

	GetRecipeDetails(recipeId int) (models.Recipe, error)

	InsertRecipe(title string, userId int) (int64, error)

	InsertIngredient(name, amount, unit string, recipeId int64) error

	InsertDirection(direction string, recipeId int64) error

	GetUserByEmail(email, password string) (models.User, error)

	InsertUser(name, email, password string) (models.User, error)

	UpdateRecipe(jsonRecipe models.JsonRecipe) (models.Recipe, error)

	DeleteRecipe(id int) error
}
