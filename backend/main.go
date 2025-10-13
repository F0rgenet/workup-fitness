package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"workup_fitness/config"
	"workup_fitness/domain/user"
	"workup_fitness/middleware"
)

func main() {
	config.LoadConfig()

	mux := http.NewServeMux()

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
	userHandler := user.NewHandler(userService, config.JwtSecret)

	mux.HandleFunc("/users/register", userHandler.Register)
	mux.HandleFunc("/users/login", userHandler.Login)
	mux.Handle("/me", middleware.Auth(config.JwtSecret)(http.HandlerFunc(userHandler.Me)))

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}