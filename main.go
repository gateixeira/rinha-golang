package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler := &Service{
		prefix: "/clientes/",
		store:  Storage{DB: db},
	}
	http.Handle("/clientes/", handler)

	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
