package server

import (
	md "./models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
}

//var DB DataBase

func (a *App) Initialize(user, password, dbname, host, port string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) CreateTables() {
	if _, err := a.DB.Exec(md.UserTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(md.ForumTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(md.ThreadTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(md.PostTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(md.VoteTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(md.ForumUsersTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}
