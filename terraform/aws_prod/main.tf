terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.27.0"
    }
  }

  backend "s3" {
    bucket = "tf-state-533267191229"
    key    = "tf-state"
    region = "ap-southeast-1"
  }
}

provider "aws" {
  # Configuration options
  region  = "ap-southeast-1"
  profile = "admin-access"
}
