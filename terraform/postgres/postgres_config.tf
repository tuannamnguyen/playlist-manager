resource "postgresql_grant" "all_privileges" {
  database    = "playlist_manager"
  role        = "playlist_manager_be"
  privileges  = ["SELECT", "INSERT", "UPDATE", "DELETE"]
  object_type = "table"
  schema      = "public"
}
