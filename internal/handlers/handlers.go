package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/popnfresh234/recipe-app-golang/internal/config"
	"github.com/popnfresh234/recipe-app-golang/internal/driver"
	"github.com/popnfresh234/recipe-app-golang/internal/forms"
	"github.com/popnfresh234/recipe-app-golang/internal/models"
	"github.com/popnfresh234/recipe-app-golang/internal/renderer"
	"github.com/popnfresh234/recipe-app-golang/repository"
	"github.com/popnfresh234/recipe-app-golang/repository/dbrepo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

// NewRepo creates a new repository with an app config
func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewMysqlRepo(db.SQL, app),
	}
}

// NewHandlers sets the repository
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the homepage handler
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {

	recipes, err := repo.DB.GetAllRecipes()
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Error getting recipes")
		return
	}
	templateData := &models.TemplateData{}
	data := make(map[string]interface{})
	data["recipes"] = recipes
	if repo.App.Session.Exists(r.Context(), "user") {
		user := repo.App.Session.Get(r.Context(), "user")
		data["user"] = user
	}
	templateData.Data = data

	_ = renderer.Template(w, r, "home.page.tmpl", templateData)
}

// RecipeDetails looks up a recipe by its ID
func (repo *Repository) RecipeDetails(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	recipeID, err := strconv.Atoi(exploded[3])
	if err != nil {
		fmt.Println(err)
		repo.App.Session.Put(r.Context(), "error", "Error getting recipe details")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	recipe, err := repo.DB.GetRecipeDetails(recipeID)
	data := make(map[string]interface{})
	data["recipe"] = recipe

	recipeJson, err := json.Marshal(recipe)
	if err != nil {
		fmt.Println(err)
		repo.App.Session.Put(r.Context(), "error", "Error parsing JSON")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	data["recipeJson"] = string(recipeJson)

	// Create template data
	td := models.TemplateData{
		Data: data,
	}

	user := repo.App.Session.Get(r.Context(), "user")
	if user != nil {
		if user.(models.User).Email == recipe.User.Email {
			td.IsAuthor = true
		}
	}
	_ = renderer.Template(w, r, "recipe-details.page.tmpl", &td)
}

// EditRecipe is the handler for editing a recipe
func (repo *Repository) EditRecipe(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	recipeID, err := strconv.Atoi(exploded[3])
	if err != nil {
		fmt.Println(err)
		repo.App.Session.Put(r.Context(), "error", "Error getting recipe details")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	recipe, err := repo.DB.GetRecipeDetails(recipeID)
	data := make(map[string]interface{})
	data["recipe"] = recipe

	recipeJson, err := json.Marshal(recipe)
	if err != nil {
		fmt.Println(err)
		repo.App.Session.Put(r.Context(), "error", "Error parsing JSON")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	data["recipeJson"] = string(recipeJson)

	// Create template data
	td := models.TemplateData{
		Data: data,
	}

	user := repo.App.Session.Get(r.Context(), "user")
	if user != nil {
		if user.(models.User).Email == recipe.User.Email {
			td.IsAuthor = true
		}
	}
	_ = renderer.Template(w, r, "recipe-edit.page.tmpl", &td)

}

// PostEditRecipe handles posting an edited recipe
func (repo *Repository) PostEditRecipe(w http.ResponseWriter, r *http.Request) {

	var updatedRecipe models.JsonRecipe
	err := json.NewDecoder(r.Body).Decode(&updatedRecipe)
	if err != nil {
		fmt.Println(err)
		// TODO handle Error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = repo.DB.UpdateRecipe(updatedRecipe)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Insert ingredients
	for _, ingredient := range updatedRecipe.Ingredients {
		err = repo.DB.InsertIngredient(ingredient.Name, ingredient.Amount, ingredient.Unit, int64(updatedRecipe.ID))
		if err != nil {
			//TODO Delete recipe?
			fmt.Println(err)
		}
	}

	// Insert directions
	for _, direction := range updatedRecipe.Directions {
		err = repo.DB.InsertDirection(direction.Direction, int64(updatedRecipe.ID))
		if err != nil {
			//TODO Delete recipe?
			fmt.Println(err)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Login is the login handler
func (repo *Repository) Login(w http.ResponseWriter, r *http.Request) {
	if repo.App.Session.Exists(r.Context(), "user") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	_ = renderer.Template(w, r, "login.page.tmpl", &models.TemplateData{Form: forms.New(nil)})
}

// PostLogin handles logging the user in
func (repo *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{}
	user.Email = r.Form.Get("email")
	user.Password = r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = user
		_ = renderer.Template(w, r, "login.page.tmpl", &models.TemplateData{Data: data, Form: form})
		return
	}

	user, err = repo.DB.GetUserByEmail(user.Email, user.Password)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Error signing in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	repo.App.Session.Put(r.Context(), "user", user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout is the logout handler
func (repo *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	repo.App.Session.Remove(r.Context(), "user")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Signup is the signup handler
func (repo *Repository) Signup(w http.ResponseWriter, r *http.Request) {
	_ = renderer.Template(w, r, "signup.page.tmpl", &models.TemplateData{Form: forms.New(nil)})
}

// PostSignup handles user signup form
func (repo *Repository) PostSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{}
	user.Email = r.Form.Get("email")
	user.Name = r.Form.Get("name")
	user.Password = r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MinLength("password", 6)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = user
		_ = renderer.Template(w, r, "signup.page.tmpl", &models.TemplateData{Data: data, Form: form})
		return
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		data := make(map[string]interface{})
		data["user"] = user
		repo.App.Session.Put(r.Context(), "error", "Bcrypt Error")
		_ = renderer.Template(w, r, "signup.page.tmpl", &models.TemplateData{Data: data, Form: form})
		return
	}
	user.Password = string(hashedPwd)
	user, err = repo.DB.InsertUser(user.Name, user.Email, user.Password)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Error inserting new user")
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
		return
	}
	repo.App.Session.Put(r.Context(), "user", user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// NewRecipe handles creating a new recipe
func (repo *Repository) NewRecipe(w http.ResponseWriter, r *http.Request) {
	_ = renderer.Template(w, r, "new-recipe.page.tmpl", &models.TemplateData{})
}

// PostNewRecipe handles posting a new recipe to the server
func (repo *Repository) PostNewRecipe(w http.ResponseWriter, r *http.Request) {
	var newRecipe models.JsonRecipe
	err := json.NewDecoder(r.Body).Decode(&newRecipe)
	if err != nil {
		fmt.Println(err)
		// TODO handle Error
		w.WriteHeader(http.StatusInternalServerError)
	}

	// TODO Validate Recipe
	// TODO Validate Ingredients
	// TODO Validate Directions
	// TODO Create and Insert Recipe, get ID
	user := repo.App.Session.Get(r.Context(), "user").(models.User)
	fmt.Println(user.ID)
	recipeId, err := repo.DB.InsertRecipe(newRecipe.Title, newRecipe.Image, user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Insert ingredients
	for _, ingredient := range newRecipe.Ingredients {
		err = repo.DB.InsertIngredient(ingredient.Name, ingredient.Amount, ingredient.Unit, recipeId)
		if err != nil {
			//TODO Delete recipe?
			fmt.Println(err)
		}
	}

	// Insert directions
	for _, direction := range newRecipe.Directions {
		err = repo.DB.InsertDirection(direction.Direction, recipeId)
		if err != nil {
			//TODO Delete recipe?
			fmt.Println(err)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PostDeleteRecipe deletes a recipe from the database
func (repo *Repository) PostDeleteRecipe(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	recipeID, err := strconv.Atoi(exploded[3])
	if err != nil {
		//TODO handle err
		fmt.Println(err)
	}

	err = repo.DB.DeleteRecipe(recipeID)
	if err != nil {
		//TODO handle err
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
