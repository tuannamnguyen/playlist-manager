resource "google_storage_bucket" "backend_test_bucket" {
  name          = "playlist-manager-test-bucket"
  location      = "ASIA-SOUTHEAST1"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_binding" "iam_binding_storage_user" {
  bucket  = google_storage_bucket.backend_test_bucket.name
  role    = "roles/storage.objectUser"
  members = ["serviceAccount:${google_service_account.backend_test_service_account.email}"]
}
