package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/natasha-m-oliveira/crud-in-memory/api"
	"github.com/natasha-m-oliveira/crud-in-memory/db"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
		os.Exit(1)
	}
	slog.Info("all systems offline")
}

func run() error {
	ur := db.NewUsersRepository()

	handler := api.NewHandler(ur)

	s := http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
