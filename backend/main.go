package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"workup_fitness/domain/user"
)

func main() {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := user.NewSQLiteRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	http.HandleFunc("/users", userHandler.Register)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}