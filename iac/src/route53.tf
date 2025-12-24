locals {
  vps_domain = var.environment == "prod" ? "vps.${data.aws_route53_zone.main.name}" : "vps.${var.environment}.${data.aws_route53_zone.main.name}"
}

data "aws_route53_zone" "main" {
  zone_id = var.route53_zone_id
}

resource "aws_route53_record" "vps" {
  zone_id = var.route53_zone_id
  name    = local.vps_domain
  type    = "A"
  ttl     = 300
  records = [var.vps_ip]
}