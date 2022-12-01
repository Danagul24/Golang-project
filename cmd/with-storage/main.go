package main

import (
	"context"
	"log"
	"project/internal/http"
	"project/internal/store/postgres"
)

func main() {
	urlExample := "postgres://postgres:dana@localhost:5432/postgres"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}
	srv.WaitForGracefulTermination()
}
