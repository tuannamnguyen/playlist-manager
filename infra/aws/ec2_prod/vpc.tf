module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "my-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway   = false
  enable_vpn_gateway   = false
  enable_dns_support   = true
  enable_dns_hostnames = true


  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

resource "aws_security_group" "endpoint-sg" {
  name        = "endpoint_access"
  description = "Allow inbound traffic"
  vpc_id      = module.vpc.vpc_id
}

resource "aws_vpc_security_group_ingress_rule" "ssm_https_traffic" {
  security_group_id = aws_security_group.endpoint-sg.id
  ip_protocol       = "tcp"
  cidr_ipv4         = module.vpc.vpc_cidr_block
  from_port         = "443"
  to_port           = "443"
  description       = "Enable access for the AWS SSM endpoints"
}

# Define all endpoints in a map
locals {
  interface_endpoints = {
    # Required endpoints for Systems Manager
    # https://docs.aws.amazon.com/systems-manager/latest/userguide/setup-create-vpc.html#create-vpc-endpoints

    ssm         = "com.amazonaws.${data.aws_region.current.region}.ssm"
    ssmmessages = "com.amazonaws.${data.aws_region.current.region}.ssmmessages"
    ec2messages = "com.amazonaws.${data.aws_region.current.region}.ec2messages"
  }
}

# Create all Interface Endpoints dynamically
resource "aws_vpc_endpoint" "interface" {
  for_each = local.interface_endpoints

  vpc_id              = module.vpc.vpc_id
  service_name        = each.value
  vpc_endpoint_type   = "Interface"
  subnet_ids          = module.vpc.private_subnets
  security_group_ids  = [aws_security_group.endpoint-sg.id]
  private_dns_enabled = true

  tags = {
    Name = "${each.key}-endpoint"
  }
}

# add this so EC2 instance can access S3 Ansible bucket
# https://docs.ansible.com/projects/ansible/latest/collections/amazon/aws/aws_ssm_connection.html#requirements
resource "aws_vpc_endpoint" "s3_gateway_endpoint" {
  vpc_id            = module.vpc.vpc_id
  service_name      = "com.amazonaws.${data.aws_region.current.region}.s3"
  vpc_endpoint_type = "Gateway"

  route_table_ids = module.vpc.private_route_table_ids

  tags = {
    Name = "s3-gw-endpoint"
  }
}
