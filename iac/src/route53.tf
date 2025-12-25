locals {
  vps_domain     = var.environment == "prod" ? "vps.${data.aws_route53_zone.main.name}" : "vps.${var.environment}.${data.aws_route53_zone.main.name}"
  backend_domain = var.environment == "prod" ? "q-backend.${data.aws_route53_zone.main.name}" : "q-backend.${var.environment}.${data.aws_route53_zone.main.name}"
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

resource "aws_route53_record" "backend" {
  zone_id = var.route53_zone_id
  name    = local.backend_domain
  type    = "A"
  ttl     = 300
  records = [aws_eip.backend.public_ip]
}