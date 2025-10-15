package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"

	"workup_fitness/config"
	"workup_fitness/domain/auth"
	"workup_fitness/domain/user"
)

func main() {
	config.LoadConfig()

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

	authService := auth.NewService(*userService)
	authHandler := auth.NewHandler(authService, config.JwtSecret)

	r := chi.NewRouter()
	user.RegisterRoutes(r, userHandler)
	auth.RegisterRoutes(r, authHandler)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
