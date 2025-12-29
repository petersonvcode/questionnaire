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
    server_assets_bucket           = aws_s3_bucket.server_assets.bucket
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

resource "aws_ssm_parameter" "backend_ssh_pub_key" {
  name = "q-backend-ssh-pub-key-${var.environment}"
  type = "String"

  value = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN5ToLhK3EecKfkJkZOjb+9n2AGD44IiEuqHWjavYE5H q-user@questionnaire.com"
}