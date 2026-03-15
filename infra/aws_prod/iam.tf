

resource "aws_iam_policy" "ec2_app_instance_policy" {
  name        = "playlist-manager-ec2-policy"
  description = "Policy to provide EC2 permission to have access to S3"
  policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = [
        {
          Effect = "Allow",
          Action = [
            "s3:PutObject",
            "s3:GetObject"
          ]
          Resource = "${aws_s3_bucket.example.arn}/*"
        },
      ]

    }
  )
}

resource "aws_iam_role" "ec2_app_instance_role" {
  name = "playlist-manager-ec2-role"
  assume_role_policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = [
        {
          Action = "sts:AssumeRole"
          Effect = "Allow"
          Principal = {
            Service = "ec2.amazonaws.com"
          }
        }
      ]
    }
  )
}

resource "aws_iam_role_policy_attachment" "iam_role_attachment" {
  role       = aws_iam_role.ec2_app_instance_role.name
  policy_arn = aws_iam_policy.ec2_app_instance_policy.arn
}

resource "aws_iam_instance_profile" "ec2_profile" {
  name = "playlist-manager-profile"
  role = aws_iam_role.ec2_app_instance_role.name
}
