# VPC Network
resource "google_compute_network" "backend_vpc_network" {
  name                    = "playlist-manager-vpc-network"
  auto_create_subnetworks = false
}
