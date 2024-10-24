variable "db_root_password" {
  type = string
}

data "google_project" "project" {
}

resource "google_sql_database_instance" "prod_db" {
  name             = "playlist-manager-prod"
  region           = local.region
  database_version = "POSTGRES_16"
  root_password    = var.db_root_password

  settings {
    tier            = "db-g1-small"
    disk_autoresize = false
  }
  deletion_protection = false
  depends_on          = [google_project_service.gcp_services]
}

resource "google_sql_database" "playlist_manager_database" {
  name     = "playlist_manager"
  instance = google_sql_database_instance.prod_db.name
}

resource "random_password" "db_user_pwd" {
  length  = 16
  special = false
}

resource "google_sql_user" "google_sql_user" {
  name     = "user"
  instance = google_sql_database_instance.prod_db.name
  password = random_password.db_user_pwd.result
}
