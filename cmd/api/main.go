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
	ctx := context.Background()

	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, os.Getenv("PG_CON_STR"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Create sqlc queries struct
	queries := sqlc.New(pool)

	// Create validator object
	validate = validator.New()

	// Create chi router
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(20 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"SERVER": "RUNNING..."})
	})

	// DOMAINS --------- ----------- -----------
	// User domain
	userRepo := user.NewUserRepository(queries)
	userSvc := user.NewUserService(userRepo)
	user.RegisterValidations(validate)
	userHandler := user.NewUserHandler(userSvc, validate)

	// We connect base router for api/v1
	r.Route("/api/v1", func(r chi.Router) {
		// user routes
		r.Mount("/users", userHandler.Routes(r))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, http.StatusNotFound, errors.New("this route not defined"))
	})

	// Start server
	addr := os.Getenv("PORT")
	if addr == "" {
		log.Fatal("PORT env var is required (e.g. :8080)")
	}
	if addr[0] != ':' {
		addr = ":" + addr
	}

	fmt.Printf("server listening PORT %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
