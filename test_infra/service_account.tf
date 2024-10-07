resource "google_service_account" "backend_test_service_account" {
  account_id                   = "playist-manager-backend-test"
  display_name                 = "playist-manager-backend-test"
  create_ignore_already_exists = true
}
