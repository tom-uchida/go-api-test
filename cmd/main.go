// testcontainers-go を使用して Spanner Emulator を起動するサンプル

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/tom-uchida/go-spanner-emulator/internal"
)

const port = "8080"

func main() {
	ctx := context.Background()

	if container, err := internal.InitSpannerEmulator(ctx); err != nil {
		log.Fatalf("failed to start spanner emulator: %v", err)
	} else {
		defer container.Terminate(ctx)
	}

	if err := internal.CreateSpannerInstance(ctx); err != nil {
		log.Fatalf("failed to create spanner instance: %v", err)
	}

	http.HandleFunc("/drop-database", internal.DropDatabase)
	http.HandleFunc("/create-database", internal.CreateDatabase)
	http.HandleFunc("/create-user", internal.CreateUser)
	http.HandleFunc("/get-user", internal.GetUser)

	fmt.Println("")
	log.Printf("Server running at: localhost:%s", port)
	fmt.Println("")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
