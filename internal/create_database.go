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
	dbpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

type CreateDatabaseReq struct {
	DatabaseName string `json:"database_name"`
	TableName    string `json:"table_name"`
}

type CreateDatabaseRes struct {
	DatabaseName string `json:"database_name"`
	TableName    string `json:"table_name"`
}

func CreateDatabase(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req CreateDatabaseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatalf("failed to decode request body: %v", err)
	}
	databaseName := req.DatabaseName
	tableName := req.TableName
	if databaseName == "" || tableName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := createDatabase(ctx, databaseName, tableName); err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CreateDatabaseRes{
		DatabaseName: databaseName,
		TableName:    tableName,
	})
}

func createDatabase(ctx context.Context, databaseName, tableName string) error {
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

	dbOp, err := dbAdmin.CreateDatabase(ctx, &dbpb.CreateDatabaseRequest{
		Parent:          parent,
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", databaseName),
		ExtraStatements: []string{getTableDefinition(tableName)},
	})
	if err != nil {
		return err
	}
	if _, err := dbOp.Wait(ctx); err != nil {
		return err
	}
	fmt.Println("Database created:", databaseName)
	fmt.Println("Table created:", tableName)
	fmt.Println("")

	return nil
}

func getTableDefinition(tableName string) string {
	switch tableName {
	case "Users":
		return `CREATE TABLE Users (
				UserID   STRING(36) NOT NULL,
				Name     STRING(1024),
			) PRIMARY KEY(UserID)`
	default:
		return ""
	}
}
