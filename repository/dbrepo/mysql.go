package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/popnfresh234/recipe-app-golang/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func (dbRepo *mysqlDBRepo) GetAllRecipes() ([]models.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbRepo.DB.QueryContext(ctx, `
		select 
		    id, image, title, created_at, updated_at, user_id
		from recipes
	`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var recipes []models.Recipe
	for rows.Next() {
		var id, userId int
		var createdAt, updatedAt []byte
		var title string
		var image sql.RawBytes
		var imageBytes []byte

		err = rows.Scan(
			&id, &image, &title, &createdAt, &updatedAt, &userId,
		)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		if image != nil {
			imageBytes = []byte(image)
		}

		parsedCreated, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			fmt.Println("Error parsing created at")
			return nil, err
		}
		parsedUpdated, err := time.Parse("2006-01-02 15:04:05", string(updatedAt))
		if err != nil {
			fmt.Println("Error parsing updated at")
			return nil, err
		}
		recipe := models.Recipe{
			ID:        id,
			Image:     &imageBytes,
			Title:     title,
			CreatedAt: parsedCreated,
			UpdatedAt: parsedUpdated,
			UserId:    userId,
		}

		var user models.User
		userStatemet := `
			SELECT name
			FROM 
			    users
			WHERE
			    users.id = ?
		`

		userRow := dbRepo.DB.QueryRowContext(ctx, userStatemet, recipe.UserId)
		err = userRow.Scan(&user.Name)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		recipe.User = user
		recipes = append(recipes, recipe)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

// GetRecipeDetails gets a recipe by ID
func (dbRepo *mysqlDBRepo) GetRecipeDetails(recipeId int) (models.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	recipeStatement := `
		SELECT
			id, title, user_id
		FROM
		    recipes
		WHERE
		    recipes.id = ?
	`
	var recipe models.Recipe

	recipeRow := dbRepo.DB.QueryRowContext(ctx, recipeStatement, recipeId)
	err := recipeRow.Scan(&recipe.ID, &recipe.Title, &recipe.UserId)
	if err != nil {
		fmt.Println(err)
		return recipe, err
	}

	// Get User
	userStatement := `
		SELECT
			id, name, email
		FROM
		    users
		WHERE
		    users.id = ?
	`

	var user models.User
	userRow := dbRepo.DB.QueryRowContext(ctx, userStatement, recipe.UserId)
	err = userRow.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		fmt.Println(err)
		return recipe, err
	}
	recipe.User = user

	// Get ingredients
	ingStatement := `
			SELECT
		    	id, name,amount, unit
			FROM
		    	ingredients
			WHERE
		    	recipe_id = ?
		`
	var ingredients []models.Ingredient
	ingredientRows, err := dbRepo.DB.QueryContext(ctx, ingStatement, recipeId)
	for ingredientRows.Next() {
		var ingredient models.Ingredient
		err = ingredientRows.Scan(
			&ingredient.ID, &ingredient.Name, &ingredient.Amount, &ingredient.Unit,
		)
		if err != nil {
			log.Fatal("Ingredient scan error")
		}
		ingredients = append(ingredients, ingredient)
	}

	//Get directions
	dirStatement := `
			SELECT
		    	id, direction
			FROM
		    	directions
			WHERE
		    	recipe_id = ?
		`
	var directions []models.Direction
	directionRows, err := dbRepo.DB.QueryContext(ctx, dirStatement, recipeId)
	for directionRows.Next() {
		var direction models.Direction
		err = directionRows.Scan(
			&direction.ID, &direction.Direction,
		)
		if err != nil {
			log.Fatal("Direction scan error")
		}
		directions = append(directions, direction)
	}

	recipe.Ingredients = ingredients
	recipe.Directions = directions
	return recipe, nil
}

func (dbRepo *mysqlDBRepo) InsertRecipe(title string, userId int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement :=
		`INSERT INTO recipes (title,image,user_id, created_at, updated_at)
 		VALUES (?,?,?,?,?)
		`
	res, err := dbRepo.DB.ExecContext(ctx, statement, title, "", userId, time.Now(), time.Now())
	if err != nil {
		log.Println("Error inserting user", err)
		return -1, err
	}

	newId, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return newId, nil
}

// InsertIngredient handles inserting an ingredient into the database
func (dbRepo *mysqlDBRepo) InsertIngredient(name, amount, unit string, recipeId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement :=
		`INSERT INTO ingredients (name,amount,unit,recipe_id, created_at, updated_at)
 		VALUES (?,?,?,?,?,?)
		`
	_, err := dbRepo.DB.ExecContext(ctx, statement, name, amount, unit, recipeId, time.Now(), time.Now())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// InsertDirection handles inserting a direction into the database
func (dbRepo *mysqlDBRepo) InsertDirection(direction string, recipeId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement :=
		`INSERT INTO directions (direction,recipe_id, created_at, updated_at)
 		VALUES (?,?,?,?)
		`
	_, err := dbRepo.DB.ExecContext(ctx, statement, direction, recipeId, time.Now(), time.Now())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetUserByEmail looks up a user by email
func (dbRepo *mysqlDBRepo) GetUserByEmail(email, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := dbRepo.DB.QueryRowContext(ctx, `SELECT id, name, image, email, password, created_at, updated_at FROM users WHERE email = ?`, email)

	var user models.User
	var createdAt, updatedAt []byte

	err := row.Scan(&user.ID, &user.Name, &user.Image, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		log.Println("Error fetching newly created user", err)
		return models.User{}, err
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		log.Println("Error parsing time", err)
		return models.User{}, err
	}

	user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		log.Println("Error parsing time", err)
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return models.User{}, errors.New("incorrect password")
	}

	return user, nil
}

// InsertUser inserts a new user into the DB
func (dbRepo *mysqlDBRepo) InsertUser(name, email, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement :=
		`INSERT INTO users (name, image,email, password, created_at, updated_at)
 		VALUES (?,?,?,?,?,?)
		`

	res, err := dbRepo.DB.ExecContext(ctx, statement, name, "", email, password, time.Now(), time.Now())

	if err != nil {
		log.Println("Error inserting user", err)
		return models.User{}, err
	}

	newId, err := res.LastInsertId()

	if err != nil {
		log.Println("Error inserting user", err)
		return models.User{}, err
	}

	row := dbRepo.DB.QueryRowContext(ctx, `SELECT * FROM users WHERE id = ?`, newId)
	var user models.User
	var createdAt, updatedAt []byte

	err = row.Scan(&user.ID, &user.Name, &user.Image, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		log.Println("Error fetching newly created user", err)
		return models.User{}, err
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		log.Println("Error parsing time", err)
		return models.User{}, err
	}

	user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		log.Println("Error parsing time", err)
		return models.User{}, err
	}
	return user, nil
}

func (dbRepo *mysqlDBRepo) UpdateRecipe(jsonRecipe models.JsonRecipe) (models.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `UPDATE recipes SET title = ?, updated_at = ? WHERE id = ?`
	_, err := dbRepo.DB.ExecContext(ctx, statement, jsonRecipe.Title, time.Now(), jsonRecipe.ID)
	if err != nil {
		return models.Recipe{}, err
	}

	// Delete old ingredients
	deleteIngredients := `DELETE FROM ingredients WHERE recipe_id = ?`
	_, err = dbRepo.DB.ExecContext(ctx, deleteIngredients, jsonRecipe.ID)
	if err != nil {
		return models.Recipe{}, err
	}

	// Delete old directions
	deleteDirections := `DELETE FROM directions WHERE recipe_id = ?`
	_, err = dbRepo.DB.ExecContext(ctx, deleteDirections, jsonRecipe.ID)
	if err != nil {
		return models.Recipe{}, err
	}

	return models.Recipe{}, err
}

func (dbRepo *mysqlDBRepo) DeleteRecipe(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `DELETE FROM recipes WHERE id = ?`
	_, err := dbRepo.DB.ExecContext(ctx, statement, id)
	if err != nil {
		return err
	}

	return nil
}
