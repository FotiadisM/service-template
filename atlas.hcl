env "local" {
  src = "file://internal/db/schema"
  url = getenv("ATLAS_PSQL_URI")
  dev = getenv("ATLAS_PSQL_DEV_URI")

  migration {
    dir    = "file://internal/db/migrations"
    format = atlas
  }

  format {
    migrate {
      diff = "{{ sql . \"    \" }}"
    }
  }
}
