package main

import (
	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//
	sw "./go"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")

	sw.DB.Initialize("postgres", "admin", "forum-test", "localhost", "5432")
	defer sw.DB.DB.Close()

	sw.DB.CreateTables()

	router := sw.NewRouter()
	
	log.Fatal(http.ListenAndServe(":8080", router))
}
