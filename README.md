# go-spanner-emulator

## Start server

```shell
> go run cmd/main.go
Spanner emulator running at: localhost:55457
Instance created: projects/test-project/instances/test-instance
Database created: projects/test-project/instances/test-instance/databases/test-db

2025/09/14 22:30:56 Server running on port 8080
```

## API test

### Get User

```shell
> runn run runbook/get_user.yaml
{
  "user_id": "923a2d97-d55a-4c71-9145-9290c24ea296"
}
{
  "user_id": "5f49170d-5ff5-4fae-a9a0-2cd005863a9b"
}
{
  "users": [
    {
      "name": "test-name-1",
      "user_id": "923a2d97-d55a-4c71-9145-9290c24ea296"
    },
    {
      "name": "test-name-2",
      "user_id": "5f49170d-5ff5-4fae-a9a0-2cd005863a9b"
    }
  ]
}
.

1 scenario, 0 skipped, 0 failures
```

### SQL

```shell
> SPANNER_EMULATOR_HOST=localhost:55457 runn run runbook/sql.yaml
[
  {
    "table_name": "Users"
  }
]
.

1 scenario, 0 skipped, 0 failures
```
