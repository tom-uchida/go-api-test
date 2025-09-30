# go-api-test

## Start server

```shell
> go run cmd/main.go
Spanner emulator running at: localhost:60780
Instance created: projects/test-project/instances/test-instance

2025/09/21 23:26:45 Server running at: localhost:8080

```

## API test with runn

```shell
> SPANNER_EMULATOR_HOST=localhost:60780 runn run runbook/create_user.yaml
{
  "database_name": "create-user",
  "table_name": "Users"
}
{
  "user_id": "917f3ba9-c4f2-4f18-a8ea-d36b03be7e74"
}
{
  "user_id": "2ed531a5-9d98-414f-9136-1a77a8e9512b"
}
[
  {
    "Name": "test-name-1",
    "UserID": "917f3ba9-c4f2-4f18-a8ea-d36b03be7e74"
  },
  {
    "Name": "test-name-2",
    "UserID": "2ed531a5-9d98-414f-9136-1a77a8e9512b"
  }
]
{
  "database_name": "create-user"
}
.

1 scenario, 0 skipped, 0 failures
```

## Go test

```shell
go test -v -run TestSomething ./test
2025/09/30 14:52:17 github.com/testcontainers/testcontainers-go - Connected to docker: 
  Server Version: 28.2.2
  API Version: 1.50
  Operating System: Docker Desktop
  Total Memory: 7836 MB
  Labels:
    com.docker.desktop.address=unix:///Users/uchidatomomasa/Library/Containers/com.docker.docker/Data/docker-cli.sock
  Testcontainers for Go Version: v0.38.0
  Resolved Docker Host: unix:///var/run/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: f590d4a35c00ef98bc6c263da3f8aebf3e493396f58df526837771f8eec00f17
  Test ProcessID: beed3297-be16-4ac9-82d1-f724b54bd655
2025/09/30 14:52:17 🐳 Creating container for image gcr.io/cloud-spanner-emulator/emulator:latest
2025/09/30 14:52:17 🐳 Creating container for image testcontainers/ryuk:0.12.0
2025/09/30 14:52:17 ✅ Container created: beae3315da83
2025/09/30 14:52:17 🐳 Starting container: beae3315da83
2025/09/30 14:52:17 ✅ Container started: beae3315da83
2025/09/30 14:52:17 ⏳ Waiting for container id beae3315da83 image: testcontainers/ryuk:0.12.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms skipInternalCheck:false skipExternalCheck:false}
2025/09/30 14:52:17 🔔 Container is ready: beae3315da83
2025/09/30 14:52:17 ✅ Container created: 3e88228104c5
2025/09/30 14:52:17 🐳 Starting container: 3e88228104c5
2025/09/30 14:52:17 ✅ Container started: 3e88228104c5
2025/09/30 14:52:17 ⏳ Waiting for container id 3e88228104c5 image: gcr.io/cloud-spanner-emulator/emulator:latest. Waiting for: &{timeout:<nil> Log:Cloud Spanner emulator running IsRegexp:false Occurrence:1 PollInterval:100ms check:<nil> submatchCallback:<nil> re:<nil> log:[]}
2025/09/30 14:52:18 🔔 Container is ready: 3e88228104c5
=== RUN   TestSomething
2025/09/30 14:52:18 Database created: test-db
2025/09/30 14:52:18 
2025/09/30 14:52:18 user: ea1257cf-021f-4e78-ace7-c9e22446618e user-name
--- PASS: TestSomething (0.02s)
PASS
2025/09/30 14:52:18 🐳 Stopping container: 3e88228104c5
2025/09/30 14:52:19 ✅ Container stopped: 3e88228104c5
2025/09/30 14:52:19 🐳 Terminating container: 3e88228104c5
2025/09/30 14:52:19 🚫 Container terminated: 3e88228104c5
ok  	github.com/tom-uchida/go-spanner-emulator/test	2.018s
```
