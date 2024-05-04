env "local" {
  src = "file://internal/store/schema"
  url = "postgres://postgres:postgres@localhost:5432/local?sslmode=disable"
  dev = "postgres://postgres:postgres@localhost:5433/local?sslmode=disable"

  migration {
    dir    = "file://internal/store/migrations"
    format = atlas
  }

  format {
    migrate {
      diff = "{{ sql . \"    \" }}"
    }
  }
}
