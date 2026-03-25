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
  subnet_id                   = module.vpc.private_subnets[0]
  associate_public_ip_address = false
  iam_instance_profile        = aws_iam_instance_profile.ec2_profile.name

  root_block_device = {
    delete_on_termination = true
    size                  = 20
  }

  security_group_ingress_rules = {
    "8080_from_alb" : {
      "referenced_security_group_id" : aws_security_group.alb_security_group.id,
      "from_port" : 8080,
      "to_port" : 8080
    }
  }
  security_group_vpc_id = module.vpc.vpc_id

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
