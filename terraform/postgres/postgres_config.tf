resource "postgresql_grant" "table_privileges" {
  database = "playlist_manager"
  role     = "playlist_manager_be"
  privileges = [
    "SELECT",
    "INSERT",
    "UPDATE",
    "DELETE"
  ]
  object_type = "table"
  schema      = "public"
}

resource "postgresql_grant" "schema_privileges" {
  database = "playlist_manager"
  role     = "playlist_manager_be"
  privileges = [
    "SELECT",
    "USAGE",
    "UPDATE"
  ]
  object_type = "sequence"
  schema      = "public"
}
