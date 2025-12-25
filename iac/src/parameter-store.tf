resource "aws_ssm_parameter" "backend_configuration" {
  name = "q-backend-conf-${var.environment}"
  type = "String"

  value = jsonencode({
    backend_domain                 = local.backend_domain,
    persistence_volume_device_name = local.persistence_volume_device_name,
    persistence_mount_dir          = "/persistent"
    database_dir                   = "/persistent/db"
    app_dir                        = "/q"
  })
}
