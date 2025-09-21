# go-api-test

## Start server

```shell
> go run cmd/main.go
Spanner emulator running at: localhost:60780
Instance created: projects/test-project/instances/test-instance

2025/09/21 23:26:45 Server running at: localhost:8080

```

## API test

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
