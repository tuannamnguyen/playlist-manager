terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.4.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "3.6.3"
    }

    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "1.24.0"
    }
  }

  backend "gcs" {
    bucket = "playlist-manager-tf-state"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = local.project_id
  region  = local.region
}

provider "postgresql" {
  host     = google_sql_database_instance.prod_db.public_ip_address
  port     = 5432
  database = "postgres"
  username = "postgres_user"
  password = var.db_root_password
  sslmode  = "require"
}

locals {
  project_id = "playlist-manager-437214"
  region     = "asia-southeast1"
}
