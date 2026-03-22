resource "porkbun_domain_nameservers" "nameservers" {
  domain      = var.root_domain_name
  nameservers = aws_route53_zone.playlist_manager_zone.name_servers
}
