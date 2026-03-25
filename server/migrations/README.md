# Goose migrations

Goose migration files for the backend belong in this directory.

The local Compose stack exposes a dedicated `migrate` service in `compose.yaml`
that mounts this path and runs:

```bash
goose -dir /workspace/server/migrations postgres "$DATABASE_URL" up
```
