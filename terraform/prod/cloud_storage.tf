resource "google_storage_bucket" "backend_prod_bucket" {
  name          = "playlist-manager-prod-bucket"
  location      = "ASIA-SOUTHEAST1"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
}