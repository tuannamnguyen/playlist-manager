resource "aws_lb" "playlist_manager_load_balancer" {
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_security_group.id]
  subnets            = module.vpc.public_subnets
}

resource "aws_lb_target_group" "playlist_manager_target_group" {
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = module.vpc.vpc_id
  target_type = "instance"

  # Health check configuration
  health_check {
    enabled             = true
    path                = "/healthcheck"
    port                = "traffic-port"
    protocol            = "HTTP"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
    matcher             = "200"
  }
}

resource "aws_lb_target_group_attachment" "tg_attachement_a" {
  target_group_arn = aws_lb_target_group.playlist_manager_target_group.arn
  target_id        = module.ec2_instance.id
  port             = 8080
}

# HTTP listener - redirect all traffic to HTTPS
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.playlist_manager_load_balancer.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

# HTTPS listener - forward to the app target group
resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.playlist_manager_load_balancer.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06" # Modern TLS
  certificate_arn   = aws_acm_certificate_validation.cert_validation.certificate_arn


  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.playlist_manager_target_group.arn
  }
}

resource "aws_security_group" "alb_security_group" {
  name_prefix = "alb-"
  description = "Security group for the ALB"
  vpc_id      = module.vpc.vpc_id
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_app_https_traffic" {
  security_group_id = aws_security_group.alb_security_group.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.home_ip_address
  from_port         = "443"
  to_port           = "443"
}

resource "aws_vpc_security_group_ingress_rule" "allow_app_http_traffic" {
  security_group_id = aws_security_group.alb_security_group.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.home_ip_address
  from_port         = "80"
  to_port           = "80"
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv4" {
  security_group_id = aws_security_group.alb_security_group.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1" # semantically equivalent to all ports
}
