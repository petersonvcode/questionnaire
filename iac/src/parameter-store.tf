resource "aws_ssm_parameter" "backend_configuration" {
  name = "q-backend-conf-${var.environment}"
  type = "String"

  value = jsonencode({
    backend_domain                 = local.backend_domain,
    persistence_volume_device_name = local.persistence_volume_device_name,
    persistence_mount_dir          = "/persistent"
    database_dir                   = "/persistent/db"
    logs_dir                       = "/persistent/logs"
    app_dir                        = "/q"
  })
}

resource "aws_ssm_parameter" "github_token" {
  name = "q-github-token-${var.environment}"
  type = "String"

  value = "needs to be set manually"
  
  lifecycle {
    ignore_changes = [value]
  }
}
