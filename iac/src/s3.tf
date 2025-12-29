resource "aws_s3_bucket" "server_assets" {
    bucket = "q-server-assets-${var.environment}"
}