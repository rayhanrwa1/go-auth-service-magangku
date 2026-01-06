package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() *pgxpool.Pool {
	dsn := "postgres://authuser:authpass@localhost:5432/authdb?sslmode=disable"

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("PostgreSQL not connected")
	}

	log.Println("PostgreSQL connected")
	return db
}
