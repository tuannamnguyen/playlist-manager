# VPC Network
resource "google_compute_network" "backend_vpc_network" {
  name                    = "playlist-manager-vpc-network"
  auto_create_subnetworks = false
}

# Subnetwork
resource "google_compute_subnetwork" "subnet" {
  name          = "playlist-manager-subnet"
  ip_cidr_range = "10.0.0.0/24"
  network       = google_compute_network.backend_vpc_network.id
}

# Serverless VPC access connector
resource "google_vpc_access_connector" "vpc_connector" {
  name          = "playlist-manager-serverless-vpc-connector"
  network       = google_compute_network.backend_vpc_network.id
  ip_cidr_range = "10.8.0.0/24"
}
