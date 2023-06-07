package utils

import (
	"chat_backend/generated"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

func Database() (*pgxpool.Pool, *generated.Queries) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	queries := generated.New(pool)
	return pool, queries
}
