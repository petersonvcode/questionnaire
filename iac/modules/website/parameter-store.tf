resource "aws_ssm_parameter" "website" {
  name  = "${var.name}-${var.environment}-website"
  type  = "String"
  value = jsonencode({
    distribution_id = aws_cloudfront_distribution.website.id
    bucket_name     = aws_s3_bucket.website.bucket
    addresses       = join(";", [for address in local.full_addresses : "https://${address}"])
  })
}
