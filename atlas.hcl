env "local" {
  src = "file://internal/database/schema"
  url = getenv("ATLAS_PSQL_URI")
  dev = getenv("ATLAS_PSQL_DEV_URI")

  migration {
    dir    = "file://internal/database/migrations"
    format = atlas
  }

  format {
    migrate {
      diff = "{{ sql . \"    \" }}"
    }
  }
}
