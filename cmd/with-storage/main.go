package main

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"project/internal/http"
	"project/internal/pkg/auth"
	"project/internal/store/postgres"
)

const key = "secret"

func main() {

	urlExample := "postgres://postgres:dana@localhost:5432/postgres"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	manager, err := auth.NewManager(key)
	if err != nil {
		panic(err)
	}

	srv := http.NewServer(context.Background(),
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithCache(cache),
		http.WithTokenManager(manager))

	if err := srv.Run(); err != nil {
		log.Println(err)
	}
	srv.WaitForGracefulTermination()
}
