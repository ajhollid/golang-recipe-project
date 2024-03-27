package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/popnfresh234/recipe-app-golang/internal/config"
	"net/http"
)

var app *config.AppConfig

// NewHelpers sets up app config
func NewHelpers(a *config.AppConfig) {
	app = a
}

// IsAuthenticated checks if a user is authenticated
func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user")
	return exists

}

// GetJson gets a JSON string
func GetGson(data interface{}) string {
	jsonString, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(jsonString)
}
