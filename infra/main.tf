terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.4.0"
    }
  }

  backend "gcs" {
    bucket = "playlist-manager-tf-state"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = "playlist-manager-437214"
  region  = "asia-southeast1"
}
