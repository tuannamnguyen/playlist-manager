resource "google_service_account" "backend_test_service_account" {
  account_id                   = "playist-manager-backend-test"
  display_name                 = "playist-manager-backend-test"
  create_ignore_already_exists = true
}

resource "google_service_account_iam_binding" "iam_binding_token_creator" {
  service_account_id = google_service_account.backend_test_service_account.id
  role               = "roles/iam.serviceAccountTokenCreator"
  members            = ["serviceAccount:${google_service_account.backend_test_service_account.email}"]
}
