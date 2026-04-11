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

resource "aws_iam_policy" "s3_allow_ec2_policy" {
  name        = "s3-allow-ec2-policy"
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

resource "aws_iam_policy" "rds_allow_ec2_connect_policy" {
  name        = "rds-connect-allow-ec2-policy"
  description = "Policy to allow EC2 instance to connect to RDS using IAM"

  # https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.IAMPolicy.html
  policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = [
        {
          Effect = "Allow"
          Action = [
            "rds-db:connect"
          ]
          Resource = [
            "arn:aws:rds-db:${data.aws_region.current.region}:${data.aws_caller_identity.current.account_id}:dbuser:${aws_db_instance.db.id}/postgres"
          ]
        }
      ]
    }
  )
}

resource "aws_iam_policy" "sm_allow_ec2_get_secret_policy" {
  name        = "sm-allow-ec2-get-secret-policy"
  description = "Policy to allow EC2 instance to get secret from Secrets Manager"

  policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = [
        {
          Effect = "Allow"
          Action = [
            "secretsmanager:GetSecretValue"
          ]
          Resource = [
            aws_secretsmanager_secret.db_host.arn
          ]
        }
      ]
    }
  )
}

locals {
  policies_arn = tomap({
    s3_allow_ec2  = aws_iam_policy.s3_allow_ec2_policy.arn,
    ssm_managed   = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore",
    rds_allow_ec2 = aws_iam_policy.rds_allow_ec2_connect_policy.arn
    sm_allow_ec2  = aws_iam_policy.sm_allow_ec2_get_secret_policy.arn
  })
}

resource "aws_iam_role_policy_attachment" "iam_role_attachment" {
  for_each = local.policies_arn

  role       = aws_iam_role.ec2_app_instance_role.name
  policy_arn = each.value
}

resource "aws_iam_instance_profile" "ec2_profile" {
  name = "playlist-manager-profile"
  role = aws_iam_role.ec2_app_instance_role.name
}
