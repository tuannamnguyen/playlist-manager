resource "aws_db_subnet_group" "db_subnets" {
  name       = "playlist_manager_db_subnet"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_db_instance" "db" {
  identifier                  = "playlist-manager"
  instance_class              = "db.t3.micro"
  allocated_storage           = 5
  engine                      = "postgres"
  engine_version              = "18.3"
  username                    = "postgres"
  db_subnet_group_name        = aws_db_subnet_group.db_subnets.name
  vpc_security_group_ids      = [aws_security_group.db_access.id]
  manage_master_user_password = true
  skip_final_snapshot         = true

}

resource "aws_security_group" "db_access" {
  name        = "db_access"
  description = "Allow traffic from EC2 instance"
  vpc_id      = module.vpc.vpc_id

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_vpc_security_group_ingress_rule" "db_allow_ec2" {
  security_group_id            = aws_security_group.db_access.id
  ip_protocol                  = "tcp"
  referenced_security_group_id = module.ec2_instance.security_group_id
  from_port                    = "5432"
  to_port                      = "5432"
}
