// testcontainers-go を使用して Spanner Emulator を起動するサンプル

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/tom-uchida/go-api-test/internal"
	"github.com/tom-uchida/go-api-test/internal/db"
)

const port = "8080"

func main() {
	ctx := context.Background()

	parent, err := db.CreateSpannerInstance(ctx)
	if err != nil {
		log.Fatalf("failed to create spanner instance: %v", err)
	}
	dsn, err := db.CreateDatabase(ctx, parent)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	client, err := db.NewClient(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to create spanner client: %v", err)
	}
	defer client.Close()

	http.HandleFunc("/create-user", internal.CreateUserHandler(client))
	http.HandleFunc("/get-user", internal.GetUserHandler(client))

	log.Printf("\nServer running at: localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
