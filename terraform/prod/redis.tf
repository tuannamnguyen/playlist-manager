resource "google_redis_instance" "playlist_manager_redis" {
  name               = "playlist-manager-prod-redis"
  tier               = "BASIC"
  memory_size_gb     = "2"
  region             = local.region
  redis_version      = "REDIS_7_0"
  authorized_network = google_compute_network.backend_vpc_network.id
}

# VM to port forward Redis for local connection
resource "google_compute_instance" "redis_connection_instance" {
  name         = "redis-connect"
  machine_type = "f1-micro"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "ubuntu-minimal-2210-kinetic-amd64-v20230126"
    }
  }

  network_interface {
    network = "default"

  }
}
