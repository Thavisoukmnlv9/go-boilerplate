
> Replace `postgres:welcome` and `go-boilerplate` with your actual credentials and database name if needed.


## 🚀 Setup `migrate` CLI

### macOS (with Homebrew)

```bash
brew install golang-migrate


## 🧪 Run Migrations

### Apply latest migrations

```bash
migrate -path migrations -database "postgres://postgres:welcome@localhost:5432/go-boilerplate?sslmode=disable" up


## 🔁 Rollback Migrations

### Rollback last migration

```bash
migrate -path migrations -database "postgres://postgres:welcome@localhost:5432/go-boilerplate?sslmode=disable" down 1


### Rollback all migrations

```bash
migrate -path migrations -database "postgres://postgres:welcome@localhost:5432/go-boilerplate?sslmode=disable" down 1