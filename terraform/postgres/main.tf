terraform {
  required_providers {

    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "1.24.0"
    }
  }

  backend "gcs" {
    bucket = "postgres-config-tf-state"
    prefix = "terraform/state"
  }
}

variable "db_root_password" {
  type = string
}

variable "db_host" {
  type = string
}


provider "postgresql" {
  host     = var.db_host
  port     = 5432
  database = "postgres"
  username = "postgres"
  password = var.db_root_password
  sslmode  = "disable"
}
