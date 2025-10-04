package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
)

type CreateUserReq struct {
	Name string `json:"name"`
}

type CreateUserRes struct {
	UserID string `json:"user_id"`
}

func CreateUserHandler(client *spanner.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		var req CreateUserReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Fatalf("failed to decode request body: %v", err)
		}
		userID, err := createUser(ctx, client, req.Name)
		if err != nil {
			log.Fatalf("failed to execute query: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreateUserRes{
			UserID: userID,
		})
	}
}

func createUser(ctx context.Context, client *spanner.Client, userName string) (string, error) {
	userID := uuid.New().String()
	m := spanner.Insert("Users", []string{"UserID", "Name"}, []interface{}{userID, userName})
	_, err := client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		return "", err
	}
	return userID, nil
}
