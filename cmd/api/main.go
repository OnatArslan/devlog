package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/OnatArslan/devlog/internal/db/sqlc"
	"github.com/OnatArslan/devlog/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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

	fmt.Println(queries)

	// Create chi router
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(20 * time.Second))

	// We connect base router for api/v1
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			httpx.WriteJSON(w, http.StatusOK, map[string]string{"STATUS": "OK"})
		})

		// r.Mount("/users")
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, http.StatusNotFound, errors.New("this route not defined"))
	})

	http.ListenAndServe(os.Getenv("PORT"), r)

}
