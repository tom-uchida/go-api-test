package db

import (
	"context"
	"fmt"
	"os"

	dbadmin "cloud.google.com/go/spanner/admin/database/apiv1"
	dbpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instadmin "cloud.google.com/go/spanner/admin/instance/apiv1"
	instpb "cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"

	"google.golang.org/api/option"
)

func CreateSpannerInstance(ctx context.Context) (string, error) {
	// SPANNER_EMULATOR_HOST 環境変数から接続先を取得
	emulatorHost := os.Getenv("SPANNER_EMULATOR_HOST")
	if emulatorHost == "" {
		return "", fmt.Errorf("SPANNER_EMULATOR_HOST environment variable is not set")
	}
	fmt.Println("SPANNER_EMULATOR_HOST:", emulatorHost)

	instAdmin, err := instadmin.NewInstanceAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return "", err
	}
	defer instAdmin.Close()

	projectID := os.Getenv("PROJECT_ID")
	instanceID := os.Getenv("INSTANCE_ID")
	parent := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	instOp, err := instAdmin.CreateInstance(ctx, &instpb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instpb.Instance{
			Name:        parent,
			Config:      "emulator-config",
			DisplayName: "Test Instance",
			NodeCount:   1,
		},
	})
	if err != nil {
		return "", err
	}
	if _, err := instOp.Wait(ctx); err != nil {
		return "", err
	}
	fmt.Println("Instance created:", parent)

	return parent, nil
}

func CreateDatabase(ctx context.Context, parent string) (string, error) {
	dbAdmin, err := dbadmin.NewDatabaseAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return "", err
	}
	defer dbAdmin.Close()

	dbName := os.Getenv("DB_NAME")
	dbOp, err := dbAdmin.CreateDatabase(ctx, &dbpb.CreateDatabaseRequest{
		Parent:          parent,
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", dbName),
		ExtraStatements: []string{`CREATE TABLE Users (
				UserID   STRING(36) NOT NULL,
				Name     STRING(1024),
			) PRIMARY KEY(UserID)`},
	})
	if err != nil {
		return "", err
	}
	if _, err := dbOp.Wait(ctx); err != nil {
		return "", err
	}
	fmt.Println("Database created:", dbName)
	fmt.Println("")

	dsn := fmt.Sprintf("%s/databases/%s", parent, dbName)
	return dsn, nil
}
