package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"workup_fitness/config"
	"workup_fitness/domain/auth"
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
	authService := auth.NewService(*userService)

	authHandler := auth.NewHandler(authService, config.JwtSecret)
	userHandler := user.NewHandler(userService)

	mux.HandleFunc("/users/register", authHandler.Register)
	mux.HandleFunc("/users/login", authHandler.Login)
	mux.Handle("/me", middleware.Auth(config.JwtSecret)(http.HandlerFunc(userHandler.GetPublicProfile)))

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
