terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.4.0"
    }

    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "1.24.0"
    }
  }

  backend "gcs" {
    bucket = "playlist-manager-test-tf-state"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = "playlist-manager-437214"
  region  = "asia-southeast1"
}

variable "db_root_password" {
  type = string
}


provider "postgresql" {
  host     = "127.0.0.1"
  port     = 5432
  database = "postgres"
  username = "postgres"
  password = var.db_root_password
  sslmode  = "disable"
}
