package dbrepo

import (
	"database/sql"
	"github.com/popnfresh234/recipe-app-golang/internal/config"
	"github.com/popnfresh234/recipe-app-golang/repository"
)

type mysqlDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewMysqlRepo(conn *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &mysqlDBRepo{
		App: app,
		DB:  conn,
	}
}
