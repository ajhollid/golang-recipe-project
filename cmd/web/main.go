package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-sql-driver/mysql"
	"github.com/popnfresh234/recipe-app-golang/internal/config"
	"github.com/popnfresh234/recipe-app-golang/internal/driver"
	"github.com/popnfresh234/recipe-app-golang/internal/handlers"
	"github.com/popnfresh234/recipe-app-golang/internal/helpers"
	"github.com/popnfresh234/recipe-app-golang/internal/models"
	"github.com/popnfresh234/recipe-app-golang/internal/renderer"
	"log"
	"net/http"
	"os"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager

const portNumber = ":3400"

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}
	fmt.Println("DB HOST:", dbHost)
	run(dbHost)

	src := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	fmt.Printf("Starting application on port %s\n", portNumber)

	err := src.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed, dying...")
	}

}

func run(dbHost string) {

	//Register models for session
	gob.Register(models.User{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	app.InProduction = false
	app.UseCache = false

	// Connect to DB
	fmt.Println("Attempting DB connection...")
	dbCfg := mysql.Config{
		User:   "admin",
		Passwd: "password",
		Net:    "tcp",
		Addr:   dbHost,
		DBName: "recipe_go_db",
	}
	db, err := driver.ConnectSQL(dbCfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	repo := handlers.NewRepo(&app, db)
	renderer.NewRenderer(&app)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
}
