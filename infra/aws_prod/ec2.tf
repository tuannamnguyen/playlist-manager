variable "ssh_public_key" {
  type        = string
  description = "SSH Public key generated from local server"
}

variable "home_ip_address" {
  type        = string
  description = "Home IP Address"
}

module "ec2_instance" {
  source = "terraform-aws-modules/ec2-instance/aws"

  name = "my_host_01"

  instance_type               = "t3.small"
  ami                         = "ami-0ed0867532b47cc2c" # Ubuntu 24.04 AMI ID
  key_name                    = aws_key_pair.ec2_key_pair.key_name
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

  security_group_ingress_rules = {
    "ssh_from_home" : {
      "cidr_ipv4" : var.home_ip_address,
      "from_port" : 22,
      "to_port" : 22
    }

    "app_access_from_anywhere" : {
      "cidr_ipv4" : "0.0.0.0/0",
      "from_port" : 8080,
      "to_port" : 8080
    }
  }
  security_group_vpc_id = module.vpc.vpc_id

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

resource "aws_key_pair" "ec2_key_pair" {
  public_key = var.ssh_public_key
  key_name   = "playlist_manager"
}

output "my_host_01_public_ip" {
  value = module.ec2_instance.public_ip
}
