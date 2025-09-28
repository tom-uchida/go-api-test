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
go test -v -run TestSomething ./cmd
2025/09/28 23:54:16 github.com/testcontainers/testcontainers-go - Connected to docker: 
  Server Version: 25.0.2
  API Version: 1.44
  Operating System: Docker Desktop
  Total Memory: 7941 MB
  Testcontainers for Go Version: v0.38.0
  Resolved Docker Host: unix:///Users/uchidatomomasa/.docker/run/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: 42658028b18fbc08a7b7d3fba772f4aba49f789fe17e4d41f8b822a5e6fb129b
  Test ProcessID: a038bef8-7f26-4daf-8ce6-7f470171fae1
2025/09/28 23:54:17 🐳 Creating container for image gcr.io/cloud-spanner-emulator/emulator:latest
2025/09/28 23:54:17 🐳 Creating container for image testcontainers/ryuk:0.12.0
2025/09/28 23:54:17 ✅ Container created: 47b34b975ba0
2025/09/28 23:54:17 🐳 Starting container: 47b34b975ba0
2025/09/28 23:54:17 ✅ Container started: 47b34b975ba0
2025/09/28 23:54:17 ⏳ Waiting for container id 47b34b975ba0 image: testcontainers/ryuk:0.12.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms skipInternalCheck:false skipExternalCheck:false}
2025/09/28 23:54:17 🔔 Container is ready: 47b34b975ba0
2025/09/28 23:54:17 ✅ Container created: ce3709923726
2025/09/28 23:54:17 🐳 Starting container: ce3709923726
2025/09/28 23:54:17 ✅ Container started: ce3709923726
2025/09/28 23:54:17 ⏳ Waiting for container id ce3709923726 image: gcr.io/cloud-spanner-emulator/emulator:latest. Waiting for: &{timeout:<nil> Log:Cloud Spanner emulator running IsRegexp:false Occurrence:1 PollInterval:100ms check:<nil> submatchCallback:<nil> re:<nil> log:[]}
2025/09/28 23:54:18 🔔 Container is ready: ce3709923726
=== RUN   TestSomething
2025/09/28 23:54:18 Database created: test-db
2025/09/28 23:54:18 
2025/09/28 23:54:18 user: a9067508-de71-42f6-b03c-32d0436b0542 user-name
--- PASS: TestSomething (0.08s)
PASS
2025/09/28 23:54:18 🐳 Stopping container: ce3709923726
2025/09/28 23:54:18 ✅ Container stopped: ce3709923726
2025/09/28 23:54:18 🐳 Terminating container: ce3709923726
2025/09/28 23:54:18 🚫 Container terminated: ce3709923726
ok  github.com/tom-uchida/go-spanner-emulator/cmd (cached)
```
