variable "image_tag" {
  type    = string
  default = "latest"
}

variable "dotenv_key" {
  type = string
}

variable "postgres_user" {
  type    = string
  default = "user"
}

variable "postgres_dbname" {
  type    = string
  default = "playlist_manager"
}

resource "google_cloud_run_v2_service" "playlist_manager_cloud_run_config" {
  name                = "playlist-manager-cloudrun"
  location            = local.region
  client              = "terraform"
  deletion_protection = false

  template {
    service_account = google_service_account.backend_prod_service_account.email
    containers {
      image = "tuannamnguyen290602/playlist-manager:${var.image_tag}"

      env {
        name  = "DOTENV_KEY"
        value = var.dotenv_key
      }
      # Sets a environment variable for instance connection name
      env {
        name  = "POSTGRES_HOST"
        value = "/cloudsql/${google_sql_database_instance.prod_db.connection_name}"
      }
      # Sets a secret environment variable for database user secret
      env {
        name  = "POSTGRES_USER"
        value = var.postgres_user
      }
      # Sets a secret environment variable for database password secret
      env {
        name = "POSTGRES_PASSWORD"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.dbpass.secret_id # secret name
            version = "latest"                                      # secret version number or 'latest'
          }
        }
      }
      # Sets a secret environment variable for database name secret
      env {
        name  = "POSTGRES_DBNAME"
        value = var.postgres_dbname
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      ports {
        container_port = 8080
      }

    }
    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.prod_db.connection_name]
      }
    }
  }
}

resource "google_cloud_run_v2_service_iam_member" "noauth" {
  location = google_cloud_run_v2_service.playlist_manager_cloud_run_config.location
  name     = google_cloud_run_v2_service.playlist_manager_cloud_run_config.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
