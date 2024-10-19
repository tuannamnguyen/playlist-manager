resource "google_service_account" "backend_prod_service_account" {
  account_id                   = "playist-manager-backend-prod"
  display_name                 = "playist-manager-backend-prod"
  create_ignore_already_exists = true
}


resource "google_project_iam_member" "token_creator_iam_prod" {
  project = local.project_id
  member  = "serviceAccount:${google_service_account.backend_prod_service_account.email}"
  role    = "roles/iam.serviceAccountTokenCreator"
}

resource "google_project_iam_member" "storage_object_user_iam" {
  project = local.project_id
  member  = "serviceAccount:${google_service_account.backend_prod_service_account.email}"
  role    = "roles/storage.objectUser"
}

resource "google_project_iam_member" "cloud_sql_client" {
  project = local.project_id
  member  = "serviceAccount:${google_service_account.backend_prod_service_account.email}"
  role    = "roles/cloudsql.client"
}
