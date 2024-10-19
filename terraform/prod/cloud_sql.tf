variable "db_root_password" {
  type = string
}

data "google_project" "project" {
}

resource "google_sql_database_instance" "prod_db" {
  name             = "playlist-manager-prod"
  region           = "asia-southeast1"
  database_version = "POSTGRES_16"
  root_password    = var.db_root_password

  settings {
    tier            = "db-g1-small"
    disk_autoresize = false
  }
  deletion_protection = false
  depends_on          = [google_project_service.sqladmin_api]
}
