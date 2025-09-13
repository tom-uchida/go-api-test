package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/spanner"
	dbadmin "cloud.google.com/go/spanner/admin/database/apiv1"
	dbpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instadmin "cloud.google.com/go/spanner/admin/instance/apiv1"
	instpb "cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/api/option"
)

const (
	projectID  = "test-project"
	instanceID = "test-instance"
	databaseID = "test-db"
)

func main() {
	ctx := context.Background()

	if container, err := setupSpannerEmulator(ctx); err != nil {
		log.Fatalf("failed to start spanner emulator: %v", err)
	} else {
		defer container.Terminate(ctx)
	}

	if err := createSpannerInstance(ctx); err != nil {
		log.Fatalf("failed to create spanner instance: %v", err)
	}

	dbPath, err := createSpannerDatabase(ctx)
	if err != nil {
		log.Fatalf("failed to create spanner database: %v", err)
	}

	client, err := spanner.NewClient(ctx, dbPath, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("failed to create spanner client: %v", err)
	}
	defer client.Close()

	if err := execQuery(ctx, client); err != nil {
		log.Fatalf("failed to execute query: %v", err)
	}

	time.Sleep(2 * time.Second)
}

func setupSpannerEmulator(ctx context.Context) (testcontainers.Container, error) {
	// Spanner Emulator の起動
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "gcr.io/cloud-spanner-emulator/emulator:latest",
			ExposedPorts: []string{"9010/tcp", "9020/tcp"},
			WaitingFor:   wait.ForLog("Cloud Spanner emulator running"),
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}

	// Spanner Emulator の接続先の設定
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "9010")
	emulatorHost := fmt.Sprintf("%s:%s", host, port.Port())
	os.Setenv("SPANNER_EMULATOR_HOST", emulatorHost)
	fmt.Println("Spanner emulator running at:", emulatorHost)

	return container, nil
}

func createSpannerInstance(ctx context.Context) error {
	instAdmin, err := instadmin.NewInstanceAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return err
	}
	defer instAdmin.Close()

	instPath := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)

	instOp, err := instAdmin.CreateInstance(ctx, &instpb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instpb.Instance{
			Name:        instPath,
			Config:      "emulator-config", // emulator 固有の config 名
			DisplayName: "Test Instance",
			NodeCount:   1,
		},
	})
	if err != nil {
		return err
	}
	if _, err := instOp.Wait(ctx); err != nil {
		return err
	}
	fmt.Println("Instance created:", instPath)

	return nil
}

func createSpannerDatabase(ctx context.Context) (string, error) {
	dbAdmin, err := dbadmin.NewDatabaseAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return "", err
	}
	defer dbAdmin.Close()

	parent := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	dbPath := fmt.Sprintf("%s/databases/%s", parent, databaseID)

	dbOp, err := dbAdmin.CreateDatabase(ctx, &dbpb.CreateDatabaseRequest{
		Parent:          parent,
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", databaseID),
		ExtraStatements: []string{
			`CREATE TABLE Users (
				UserID   STRING(36) NOT NULL,
				Name     STRING(1024),
			) PRIMARY KEY(UserID)`,
		},
	})
	if err != nil {
		return "", err
	}
	if _, err := dbOp.Wait(ctx); err != nil {
		return "", err
	}
	fmt.Println("Database created:", dbPath)
	fmt.Println("")

	return dbPath, nil
}

func execQuery(ctx context.Context, client *spanner.Client) error {
	m := spanner.Insert("Users", []string{"UserID", "Name"}, []interface{}{"user-id-1", "name-1"})
	_, err := client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		return err
	}

	iter := client.Single().Query(ctx, spanner.NewStatement("SELECT UserID, Name FROM Users"))
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
		fmt.Printf("User: ID=%s, Name=%s\n", id, name)
	}

	return nil
}
