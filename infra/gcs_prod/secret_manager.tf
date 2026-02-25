# Create dbpass secret
resource "google_secret_manager_secret" "dbpass" {
  secret_id = "dbpasssecret"
  replication {
    auto {}
  }
  depends_on = [google_project_service.gcp_services]
}

# Attaches secret data for dbpass secret
resource "google_secret_manager_secret_version" "dbpass_data" {
  secret      = google_secret_manager_secret.dbpass.id
  secret_data = random_password.db_user_pwd.result
}
