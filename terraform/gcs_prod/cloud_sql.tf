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
    tier                  = "db-g1-small"
    disk_autoresize       = true
    disk_autoresize_limit = 20

    ip_configuration {
      ipv4_enabled = true
      authorized_networks {
        name  = "home network"
        value = "14.191.163.98"
      }
    }

    backup_configuration {
      enabled                        = true
      location                       = local.region
      point_in_time_recovery_enabled = true
    }
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
  name     = "playlist_manager_be"
  instance = google_sql_database_instance.prod_db.name
  password = random_password.db_user_pwd.result
}
