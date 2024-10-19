resource "google_service_account" "backend_prod_service_account" {
  account_id                   = "playist-manager-backend-prod"
  display_name                 = "playist-manager-backend-prod"
  create_ignore_already_exists = true
}

locals {
  service_account_role_list = [
    "roles/iam.serviceAccountTokenCreator",
    "roles/storage.objectUser",
    "roles/cloudsql.client"
  ]
}

resource "google_project_iam_member" "service_account_roles" {
  for_each = toset(local.service_account_role_list)

  project = local.project_id
  member  = "serviceAccount:${google_service_account.backend_prod_service_account.email}"
  role    = each.key
}
