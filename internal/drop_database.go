package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/option"

	dbadmin "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

type DropDatabaseReq struct {
	DatabaseName string `json:"database_name"`
}

type DropDatabaseRes struct {
	DatabaseName string `json:"database_name"`
}

func DropDatabase(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req DropDatabaseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatalf("failed to decode request body: %v", err)
	}
	databaseName := req.DatabaseName
	if databaseName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := dropDatabase(ctx, databaseName); err != nil {
		log.Fatalf("failed to drop database: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DropDatabaseRes{
		DatabaseName: databaseName,
	})
}

func dropDatabase(ctx context.Context, databaseName string) error {
	database := fmt.Sprintf("%s/databases/%s", parent, databaseName)
	client, err := spanner.NewClient(ctx, database, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("failed to create spanner client: %v", err)
	}
	defer client.Close()

	dbAdmin, err := dbadmin.NewDatabaseAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return err
	}
	defer dbAdmin.Close()

	if err := dbAdmin.DropDatabase(ctx, &databasepb.DropDatabaseRequest{
		Database: database,
	}); err != nil {
		log.Fatalf("failed to drop database: %v", err)
	}
	fmt.Println("Database dropped:", databaseName)
	fmt.Println("")

	return nil
}
