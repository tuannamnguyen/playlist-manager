data "aws_caller_identity" "current" {}


resource "aws_s3_bucket" "example" {
  bucket        = "playlist-manager-bucket"
  force_destroy = true
}

resource "aws_s3_bucket" "ansible_ssm_bucket" {
  bucket = "ansible-ssm-bucket-${data.aws_caller_identity.current.account_id}"
}

output "ansible_ssm_bucket_name" {
  value = aws_s3_bucket.ansible_ssm_bucket.id
}
