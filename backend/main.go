package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"

	"workup_fitness/config"
	"workup_fitness/domain/auth"
	"workup_fitness/domain/user"
	"workup_fitness/pkg/logger"
)

func main() {
	logger.Setup()
	config.LoadConfig()

	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Error().Msg(err.Error())
	}

	userRepo := user.NewSQLiteRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authService := auth.NewService(userService)
	authHandler := auth.NewHandler(authService, config.JwtSecret)

	r := chi.NewRouter()
	user.RegisterRoutes(r, userHandler)
	auth.RegisterRoutes(r, authHandler)

	log.Info().Msg("Starting server on port " + config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}
