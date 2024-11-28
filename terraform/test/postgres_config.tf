resource "postgresql_grant" "all_privileges" {
  database    = "playlist_manager"
  role        = postgresql_role.playlist_manager_be.name
  privileges  = ["SELECT", "INSERT", "UPDATE", "DELETE"]
  object_type = "table"
  schema      = "public"
}

resource "postgresql_role" "playlist_manager_be" {
  name     = "playlist_manager_be"
  password = var.db_root_password
  login    = true
}
