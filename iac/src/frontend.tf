module "website" {
  source = "../modules/website"

  name                = "q"
  route53_zone_id     = var.route53_zone_id
  cost_mode           = var.cost_mode
  acm_certificate_arn = var.acm_certificate_arn
  environment         = var.environment
}