package renderer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/popnfresh234/recipe-app-golang/internal/config"
	"github.com/popnfresh234/recipe-app-golang/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var app *config.AppConfig
var pathToTemplates = "./web/templates"
var pathToComponents = "./web/templates/components"
var functions = template.FuncMap{"dict": dict, "formatTime": formatTime, "add": add}

// dict creates a dictionary of key-value pairs
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// formatTime allows formatting of time.Time in template
func formatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// add increments int a by int b
func add(a, b int) int {
	return a + b
}

// NewRenderer craetes a new renderer with access to AppConfig
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds default template data for every template
func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.Flash = app.Session.PopString(r.Context(), "flash")
	templateData.Error = app.Session.PopString(r.Context(), "error")
	templateData.Warning = app.Session.PopString(r.Context(), "warning")
	if app.Session.Exists(r.Context(), "user") {
		templateData.IsAuthenticated = 1
	}
	return templateData
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) error {
	var templateCache map[string]*template.Template
	var err error
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	}

	myTemplate, ok := templateCache[tmpl]

	if !ok {
		log.Println("could not get template from template cache")
		return errors.New("could not get template from template cache")
	}

	buffer := new(bytes.Buffer)
	data = AddDefaultData(data, r)
	err = myTemplate.Execute(buffer, data)
	if err != nil {
		return err
	}

	_, err = buffer.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	// Get all files named *.page.tmpl from the ./templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	// Range through files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, pageErr := template.New(name).Funcs(functions).ParseFiles(page)
		if pageErr != nil {
			return myCache, pageErr
		}

		// Add base layouts
		matches, matchesErr := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if matchesErr != nil {
			return myCache, matchesErr
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		// Add components
		components, err := filepath.Glob(fmt.Sprintf("%s/*.component.tmpl", pathToComponents))
		if err != nil {
			return myCache, matchesErr
		}

		if len(components) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.component.tmpl", pathToComponents))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}
	return myCache, nil
}
