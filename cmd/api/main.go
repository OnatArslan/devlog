package main

import (
	"context"
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

type config struct {
	httpAddr string
	dbConStr string
}

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, os.Getenv("PG_CON_STR"))

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	queries := sqlc.New(pool)

	fmt.Println(queries)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(20 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, "OK")
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Here we will add user post

	})

	http.ListenAndServe(os.Getenv("PORT"), r)

}
