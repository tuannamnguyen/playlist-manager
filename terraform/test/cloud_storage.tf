resource "google_storage_bucket" "backend_test_bucket" {
  name          = "playlist-manager-test-bucket"
  location      = "ASIA-SOUTHEAST1"
  storage_class = "STANDARD"

  force_destroy               = false
  uniform_bucket_level_access = true
}
