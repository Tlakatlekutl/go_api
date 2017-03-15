package server

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	md "./models"
)


type DataBase struct {
	DB *sql.DB
}

var DB DataBase

func (d *DataBase) Initialize(user, password, dbname, host, port string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port )

	var err error
	d.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *DataBase)CreateTables() {
	if _, err := d.DB.Exec(md.UserTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := d.DB.Exec(md.ForumTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := d.DB.Exec(md.ThreadTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := d.DB.Exec(md.PostTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := d.DB.Exec(md.VoteTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

