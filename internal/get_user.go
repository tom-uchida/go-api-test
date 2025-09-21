package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/option"
)

type GetUserRes struct {
	Users []*User `json:"users"`
}

type User struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	databaseName := r.Header.Get("X-Database-Name")
	if databaseName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	database := fmt.Sprintf("%s/databases/%s", parent, databaseName)

	client, err := spanner.NewClient(ctx, database, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("failed to create spanner client: %v", err)
	}
	defer client.Close()

	users, err := getUser(ctx, client)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
	}
	if users == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetUserRes{
		Users: users,
	})
}

func getUser(ctx context.Context, client *spanner.Client) ([]*User, error) {
	iter := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Users"))
	defer iter.Stop()

	var users []*User
	for {
		row, err := iter.Next()
		if err != nil {
			break
		}

		var id, name string
		if err := row.Columns(&id, &name); err != nil {
			log.Fatal(err)
		}
		users = append(users, &User{
			UserID: id,
			Name:   name,
		})
	}

	return users, nil
}
