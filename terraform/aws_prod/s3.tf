resource "aws_s3_bucket" "example" {
  bucket        = "playlist-manager-bucket"
  force_destroy = true
}
