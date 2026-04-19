
terraform {
  required_providers {
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "1.26.0"
    }
  }

  backend "s3" {
    bucket  = "playlist-manager-postgres-533267191229-ap-southeast-1-an"
    key     = "tf-state"
    region  = "ap-southeast-1"
    profile = "admin-access"
  }
}

variable "host" {
  default = "localhost"
}

provider "postgresql" {
  host             = var.host
  port             = 5432
  database         = "playlist_manager"
  username         = "nam"
  connect_timeout  = 15
  aws_rds_iam_auth = true
}
