package main

import (
	sw "./server"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")

	var App sw.App

	App.Initialize("postgres", "admin", "forum-test", "localhost", "5432")
	//App.Initialize("docker", "docker", "docker", "localhost", "5432")
	defer App.DB.Close()

	App.CreateTables()

	App.NewRouter()

	log.Fatal(http.ListenAndServe(":5000", App.Router))
}
