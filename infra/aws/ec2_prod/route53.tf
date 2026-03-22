resource "aws_route53_zone" "playlist_manager_zone" {
  name          = var.root_domain_name
  force_destroy = true
}

resource "aws_route53_record" "playlist_manager_cert_dns" {
  for_each = {
    for dvo in aws_acm_certificate.playlist_manager_ssl_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = aws_route53_zone.playlist_manager_zone.zone_id
}
