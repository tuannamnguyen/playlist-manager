variable "gcp_service_list" {
  description = "The list of apis necessary for the project"
  type        = list(string)
  default = [
    "cloudresourcemanager.googleapis.com",
    "serviceusage.googleapis.com",
    "secretmanager.googleapis.com",
    "secretmanager.googleapis.com",
    "sqladmin.googleapis.com",
    "run.googleapis.com"
  ]
}

resource "google_project_service" "gcp_services" {
  for_each           = toset(var.gcp_service_list)
  service            = each.key
  disable_on_destroy = false
}
