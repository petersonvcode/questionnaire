output "cloudfront_distribution_id" {
  value = aws_cloudfront_distribution.website.id
}

output "cloudfront_distribution_arn" {
  value = aws_cloudfront_distribution.website.arn
}

output "cloudfront_distribution_domain_name" {
  value = aws_cloudfront_distribution.website.domain_name
}

output "website_bucket_name" {
  value = aws_s3_bucket.website.bucket
}

output "website_distribution_id" {
  value = aws_cloudfront_distribution.website.id
}
