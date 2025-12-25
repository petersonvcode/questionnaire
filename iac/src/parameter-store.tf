resource "aws_ssm_parameter" "backend_configuration" {
  name = "q-backend-conf-${var.environment}"
  type = "String"
  value = jsonencode({
    backend_domain = local.backend_domain
  })
}