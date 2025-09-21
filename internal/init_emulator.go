package internal

import (
	"context"
	"fmt"
	"os"

	instadmin "cloud.google.com/go/spanner/admin/instance/apiv1"
	instpb "cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/api/option"
)

const (
	projectID  = "test-project"
	instanceID = "test-instance"
)

var (
	parent = fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
)

func InitSpannerEmulator(ctx context.Context) (testcontainers.Container, error) {
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

func CreateSpannerInstance(ctx context.Context) error {
	instAdmin, err := instadmin.NewInstanceAdminClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return err
	}
	defer instAdmin.Close()

	instOp, err := instAdmin.CreateInstance(ctx, &instpb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instpb.Instance{
			Name:        parent,
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
	fmt.Println("Instance created:", parent)

	return nil
}
