package main

import (
	sw "./server"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")

	//sw.DB.Initialize("postgres", "admin", "forum-test", "localhost", "5432")
	sw.DB.Initialize("docker", "docker", "docker", "localhost", "5432")
	defer sw.DB.DB.Close()

	sw.DB.CreateTables()

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":5000", router))
}
