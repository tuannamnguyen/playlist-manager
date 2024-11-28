resource "postgresql_grant" "all_privileges" {
  database    = "playlist_manager"
  role        = google_sql_user.google_sql_user.name
  privileges  = ["SELECT", "INSERT", "UPDATE", "DELETE"]
  object_type = "table"
  schema      = "public"
}
