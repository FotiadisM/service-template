env "local" {
  src = "file://internal/store/db_schema.sql"
  url = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
  dev = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

  migration {
    dir = "file://internal/store/migrations"
    format = atlas
  }

  format {
    migrate {
	  diff = "{{ sql . \"  \" }}"
	}
  }
}
