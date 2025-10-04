package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/spanner"
)

type GetUserRes struct {
	Users []*User `json:"users"`
}

type User struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func GetUserHandler(client *spanner.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

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
