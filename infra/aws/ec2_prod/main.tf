variable "porkbun_api_key" {
  type = string
}

variable "porkbun_secret_key" {
  type = string
}

data "aws_region" "current" {}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.28.0"
    }

    porkbun = {
      source  = "kyswtn/porkbun"
      version = "0.1.3"
    }
  }

  backend "s3" {
    bucket  = "tf-state-533267191229"
    key     = "tf-state"
    region  = "ap-southeast-1"
    profile = "admin-access"
  }
}

provider "aws" {
  # Configuration options
  region  = "ap-southeast-1"
  profile = "admin-access"
}

provider "porkbun" {
  api_key        = var.porkbun_api_key
  secret_api_key = var.porkbun_secret_key
}
