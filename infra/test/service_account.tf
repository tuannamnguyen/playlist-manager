resource "google_service_account" "backend_test_service_account" {
  account_id                   = "playist-manager-backend-test"
  display_name                 = "playist-manager-backend-test"
  create_ignore_already_exists = true
}


resource "google_project_iam_member" "token_creator_iam_test" {
  project = "playlist-manager-437214"
  member  = "serviceAccount:${google_service_account.backend_test_service_account.email}"
  role    = "roles/iam.serviceAccountTokenCreator"
}

resource "google_project_iam_member" "storage_object_user_iam" {
  project = "playlist-manager-437214"
  member  = "serviceAccount:${google_service_account.backend_test_service_account.email}"
  role    = "roles/storage.objectUser"
}
