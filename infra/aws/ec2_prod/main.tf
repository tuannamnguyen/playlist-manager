variable "porkbun_api_key" {
  type = string
}

variable "porkbun_secret_key" {
  type = string
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.28.0"
    }

    porkbun = {
      source  = "jianyuan/porkbun"
      version = "0.2.1"
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
  api_key    = var.porkbun_api_key
  secret_key = var.porkbun_secret_key
}
