env "local" {
  src = "./internal/store/schema.sql"

  url = "postgres://local_user:local_pass@localhost:5432/local?sslmode=disable"

  dev = "postgres://local_user:local_pass@localhost:5432/local?sslmode=disable"

  schemas = ["public"]
}
