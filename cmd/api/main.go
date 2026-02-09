package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/OnatArslan/devlog/internal/httpx"
	"github.com/OnatArslan/devlog/internal/sqlc"
	"github.com/OnatArslan/devlog/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var validate *validator.Validate

func main() {
	// Create a root context used during app bootstrapping.
	ctx := context.Background()

	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// Initialize PostgreSQL connection pool from environment configuration.
	pool, err := pgxpool.New(ctx, os.Getenv("PG_CON_STR"))
	if err != nil {
		log.Fatal(err)
	}
	// Ensure database resources are released on process shutdown.
	defer pool.Close()

	// Create sqlc queries struct
	queries := sqlc.New(pool)

	// Create validator object
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Create chi router
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(20 * time.Second))

	// DOMAINS --------- ----------- -----------
	// User domain
	// Wire repository, service, validations, and HTTP handlers for user module.
	userRepo := user.NewUserRepository(queries)
	userSvc := user.NewUserService(userRepo)
	user.RegisterValidations(validate)
	userHandler := user.NewUserHandler(userSvc, validate)

	// We connect base router for api/v1
	r.Route("/api/v1", func(r chi.Router) {
		// Expose a simple health endpoint for liveness checks.
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {

			httpx.WriteJSON(w, http.StatusOK, map[string]any{
				"status": "ok",
			})
		})

		// Mount user-related endpoints under /api/v1/users.
		r.Mount("/users", userHandler.Routes(chi.NewRouter()))
	})

	// Return consistent JSON error for undefined routes.
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, http.StatusNotFound, errors.New("this route not defined"))
	})

	// Start server
	addr := os.Getenv("PORT")
	if addr == "" {
		log.Fatal("PORT env var is required (e.g. :8080)")
	}
	// Normalize bare port values into :port format accepted by ListenAndServe.
	if addr[0] != ':' {
		addr = ":" + addr
	}

	// Start the HTTP server and terminate on fatal listen errors.
	fmt.Printf("server listening PORT %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
