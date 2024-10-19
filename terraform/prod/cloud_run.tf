variable "image_tag" {
  type    = string
  default = "latest"
}

variable "dotenv_key" {
  type = string
}

resource "google_cloud_run_v2_service" "playlist_manager_cloud_run_config" {
  name                = "playlist-manager-cloudrun"
  location            = local.region
  client              = "terraform"
  deletion_protection = false

  template {
    containers {
      image = "tuannamnguyen290602/playlist-manager:${var.image_tag}"

      env {
        name  = "DOTENV_KEY"
        value = var.dotenv_key
      }
      # Sets a environment variable for instance connection name
      env {
        name  = "POSTGRES_HOST"
        value = "/cloudsql/${local.project_id}:${local.region}:${google_sql_database_instance.prod_db.connection_name}"
      }
      # Sets a secret environment variable for database user secret
      env {
        name = "POSTGRES_USER"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.dbuser.secret_id # secret name
            version = "latest"                                      # secret version number or 'latest'
          }
        }
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
        name = "DB_NAME"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.dbname.secret_id # secret name
            version = "latest"                                      # secret version number or 'latest'
          }
        }
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
