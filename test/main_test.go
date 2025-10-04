package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/spanner"
	dbpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"github.com/google/uuid"
	"google.golang.org/api/option"

	dbadmin "cloud.google.com/go/spanner/admin/database/apiv1"
)

const (
	projectID  = "test-project"
	instanceID = "test-instance"
)

var (
	parent = fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
)

func TestMain(m *testing.M) {
	// SPANNER_EMULATOR_HOST 環境変数が設定されていることを確認
	emulatorHost := os.Getenv("SPANNER_EMULATOR_HOST")
	if emulatorHost == "" {
		log.Fatal("SPANNER_EMULATOR_HOST environment variable is not set")
	}

	log.Println("Using Spanner Emulator at:", emulatorHost)

	code := m.Run()
	os.Exit(code)
}

// 各テストで使う DB を作成するヘルパー
func setupDatabase(ctx context.Context, ddl string) (string, func()) {
	// データベース作成クライアント
	dbAdmin, err := dbadmin.NewDatabaseAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("failed to create db client: %v", err)
	}
	defer dbAdmin.Close()

	dbName := "db-" + uuid.New().String()
	dbOp, err := dbAdmin.CreateDatabase(ctx, &dbpb.CreateDatabaseRequest{
		Parent:          parent,
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", dbName),
		ExtraStatements: []string{ddl},
	})
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	if _, err := dbOp.Wait(ctx); err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	log.Println("Database created:", dbName)
	log.Println("")

	dsn := fmt.Sprintf("%s/databases/%s", parent, dbName)

	cleanup := func() {
		_ = dbAdmin.DropDatabase(ctx, &dbpb.DropDatabaseRequest{
			Database: dsn,
		})
	}

	return dsn, cleanup
}

func TestSomething(t *testing.T) {
	ctx := context.Background()

	// DB を作成
	ddl := `CREATE TABLE Users (
				UserID   STRING(36) NOT NULL,
				Name     STRING(1024),
			) PRIMARY KEY(UserID)`
	dsn, cleanup := setupDatabase(ctx, ddl)
	defer cleanup()

	// Spanner クライアントを作成
	client, err := spanner.NewClient(ctx, dsn, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("failed to create spanner client: %v", err)
	}
	defer client.Close()

	// ユーザー作成
	userID := uuid.New().String()
	userName := "user-name"
	m := spanner.Insert("Users", []string{"UserID", "Name"}, []interface{}{userID, userName})
	_, err = client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	// ユーザー取得
	iter := client.Single().Query(ctx, spanner.NewStatement("SELECT * FROM Users"))
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err != nil {
			break
		}

		var id, name string
		if err := row.Columns(&id, &name); err != nil {
			log.Fatal(err)
		}
		log.Printf("user: %s %s\n", id, name)
	}
}
