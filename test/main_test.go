package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/spanner"
	dbpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instadmin "cloud.google.com/go/spanner/admin/instance/apiv1"
	instpb "cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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
	ctx := context.Background()

	// Spanner Emulator を起動
	container, err := initSpannerEmulator(ctx)
	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}

	// インスタンスを作成
	if err := setupInstance(ctx); err != nil {
		log.Fatalf("failed to create instance: %v", err)
	}

	// テストを実行
	code := m.Run()

	// Spanner Emulator を停止
	if err := container.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	os.Exit(code)
}

func initSpannerEmulator(ctx context.Context) (testcontainers.Container, error) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "gcr.io/cloud-spanner-emulator/emulator:latest",
			ExposedPorts: []string{"9010/tcp", "9020/tcp"},
			WaitingFor:   wait.ForLog("Cloud Spanner emulator running"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}

	// Spanner Emulator の接続先の設定
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "9010")
	emulatorHost := fmt.Sprintf("%s:%s", host, port.Port())

	// Spanner Emulator 用の環境変数をセット
	os.Setenv("SPANNER_EMULATOR_HOST", emulatorHost)

	return container, nil
}

func setupInstance(ctx context.Context) error {
	// インスタンス作成クライアント
	instClient, err := instadmin.NewInstanceAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("failed to create instance client: %v", err)
	}
	defer instClient.Close()

	// インスタンス作成（初回のみ idempotent）
	_, _ = instClient.CreateInstance(ctx, &instpb.CreateInstanceRequest{
		Parent:     parent,
		InstanceId: instanceID,
		Instance: &instpb.Instance{
			Name:        parent,
			Config:      "emulator-config",
			DisplayName: "Test Instance",
			NodeCount:   1,
		},
	})

	return nil
}

// 各テストで使う DB を作成するヘルパー
func setupDatabase(ctx context.Context, databaseName, ddl string) (string, func()) {
	// データベース作成クライアント
	dbAdmin, err := dbadmin.NewDatabaseAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("failed to create db client: %v", err)
	}
	defer dbAdmin.Close()

	dbOp, err := dbAdmin.CreateDatabase(ctx, &dbpb.CreateDatabaseRequest{
		Parent:          parent,
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", databaseName),
		ExtraStatements: []string{ddl},
	})
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	if _, err := dbOp.Wait(ctx); err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	log.Println("Database created:", databaseName)
	log.Println("")

	dsn := fmt.Sprintf("%s/databases/%s", parent, databaseName)

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
	databaseName := "test-db"
	ddl := `CREATE TABLE Users (
				UserID   STRING(36) NOT NULL,
				Name     STRING(1024),
			) PRIMARY KEY(UserID)`
	dsn, cleanup := setupDatabase(ctx, databaseName, ddl)
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
