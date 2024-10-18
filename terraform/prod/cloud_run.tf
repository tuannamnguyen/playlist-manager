variable "image_tag" {
  type    = string
  default = "latest"
}

variable "dotenv_key" {
  type = string
}

resource "google_cloud_run_v2_service" "playlist_manager_cloud_run_config" {
  name                = "playlist-manager-cloudrun"
  location            = "asia-southeast1"
  client              = "terraform"
  deletion_protection = false

  template {
    containers {
      image = "tuannamnguyen290602/playlist-manager:${var.image_tag}"
      env {
        name  = "DOTENV_KEY"
        value = var.dotenv_key
      }
      ports {
        container_port = 8080
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
