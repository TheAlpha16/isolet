data "external_schema" "gorm" {
  program = [
    "env",
    "GOWORK=off",
    "go",
    "run",
    "./cmd/migrate",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/17/dev?search_path=public"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "local" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/17/dev?search_path=public"
  migration {
    dir = "file://migrations"
  }
  url = "postgres://postgres:postgres@localhost:5432/isolet?sslmode=disable"
}
