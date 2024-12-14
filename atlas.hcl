env "local" {
  src = "file://internal/store/schema"
  url = getenv("ATLAS_PSQL_URI")
  dev = getenv("ATLAS_PSQL_DEV_URI")

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
