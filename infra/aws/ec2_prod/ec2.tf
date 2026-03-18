variable "home_ip_address" {
  type        = string
  description = "Home IP Address"
}

module "ec2_instance" {
  source = "terraform-aws-modules/ec2-instance/aws"

  name = "my_host_01"

  instance_type               = "t3.small"
  ami                         = "ami-0ed0867532b47cc2c" # Ubuntu 24.04 AMI ID
  monitoring                  = true
  subnet_id                   = module.vpc.public_subnets[0]
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.ec2_profile.name

  root_block_device = {
    delete_on_termination = true
    size                  = 20
  }

  metadata_options = {
    "http_put_response_hop_limit" : 2,
  }

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

output "instance_id" {
  value = module.ec2_instance.id
}
