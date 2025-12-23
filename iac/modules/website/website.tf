locals {
  s3_origin_id   = "myS3Origin"
  subdomains     = [for sub in var.subdomains : var.environment == "prod" ? sub : "${sub}.${var.environment}"]
  name           = var.environment == "prod" ? var.name : "${var.name}.${var.environment}"
  addresses      = concat([local.name], local.subdomains)
  full_addresses = [for address in local.addresses : "${address}.${data.aws_route53_zone.selected.name}"]
}

#############################################
### Frontend S3 bucket for object storage ###
#############################################
resource "aws_s3_bucket" "website" {
  bucket = "${var.name}-questionnaire-${var.environment}"
}

resource "aws_s3_bucket_public_access_block" "website" {
  bucket            = aws_s3_bucket.website.id
  block_public_acls = true
}

resource "aws_s3_bucket_ownership_controls" "website" {
  bucket = aws_s3_bucket.website.bucket

  rule {
    object_ownership = "ObjectWriter"
  }
}

resource "aws_s3_bucket_acl" "website" {
  bucket = aws_s3_bucket.website.id
  acl    = "private"

  depends_on = [aws_s3_bucket_ownership_controls.website]
}

resource "aws_s3_bucket_policy" "website" {
  bucket = aws_s3_bucket.website.bucket
  policy = data.aws_iam_policy_document.website.json
}

data "aws_iam_policy_document" "website" {
  statement {
    sid       = "AllowCloudFrontServicePrincipal"
    effect    = "Allow"
    actions   = ["s3:GetObject"]
    resources = ["arn:aws:s3:::${aws_s3_bucket.website.bucket}/*"]
    principals {
      type        = "Service"
      identifiers = ["cloudfront.amazonaws.com"]
    }

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceArn"
      values   = [aws_cloudfront_distribution.website.arn]
    }
  }
}

#######################################
## Cloudfront distribution as a CDN ###
#######################################
resource "aws_cloudfront_origin_access_control" "website" {
  name                              = "${var.name}-acl-${var.environment}"
  description                       = "Website policy for ${var.name} in ${var.environment}"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "website" {
  origin {
    origin_id                = local.s3_origin_id
    origin_access_control_id = aws_cloudfront_origin_access_control.website.id
    domain_name              = aws_s3_bucket.website.bucket_regional_domain_name
  }

  aliases = local.full_addresses

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  custom_error_response {
    error_code            = 403
    error_caching_min_ttl = 10
    response_page_path    = "/index.html"
    response_code         = 200
  }

  custom_error_response {
    error_code            = 400
    error_caching_min_ttl = 10
    response_page_path    = "/index.html"
    response_code         = 200
  }

  custom_error_response {
    error_code            = 404
    error_caching_min_ttl = 10
    response_page_path    = "/index.html"
    response_code         = 200
  }

  price_class = var.cost_mode == "cheap" ? "PriceClass_100" : "PriceClass_All"

  viewer_certificate {
    acm_certificate_arn = var.acm_certificate_arn
    ssl_support_method  = "sni-only"
  }

  # Add CORS headers to allow any origin
  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = local.s3_origin_id
    compress               = true
    viewer_protocol_policy = "redirect-to-https"

    cache_policy_id            = var.enable_cors ? aws_cloudfront_cache_policy.cors_cache_policy[0].id : "658327ea-f89d-4fab-a63d-7e88639e58f6" # Managed-CachingOptimized
    response_headers_policy_id = var.enable_cors ? aws_cloudfront_response_headers_policy.cors_allow_all[0].id : null
  }

}

resource "aws_cloudfront_cache_policy" "cors_cache_policy" {
  count = var.enable_cors ? 1 : 0
  name  = "${var.name}-cors-cache-policy-${var.environment}"

  comment     = "Cache policy for CORS enabled website"
  default_ttl = 86400    # 1 day
  max_ttl     = 31536000 # 1 year
  min_ttl     = 0

  parameters_in_cache_key_and_forwarded_to_origin {
    enable_accept_encoding_brotli = true
    enable_accept_encoding_gzip   = true

    query_strings_config {
      query_string_behavior = "none"
    }

    headers_config {
      header_behavior = "whitelist"
      headers {
        items = ["Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"]
      }
    }

    cookies_config {
      cookie_behavior = "none"
    }
  }
}

resource "aws_cloudfront_response_headers_policy" "cors_allow_all" {
  count = var.enable_cors ? 1 : 0
  name  = "${var.name}-cors-allow-all-${var.environment}"

  cors_config {
    access_control_allow_origins {
      items = ["*"]
    }
    access_control_allow_headers {
      items = ["*"]
    }
    access_control_allow_methods {
      items = ["GET", "HEAD", "OPTIONS"]
    }
    origin_override                  = true
    access_control_allow_credentials = false
  }

  security_headers_config {
    content_type_options {
      override = false
    }
    frame_options {
      frame_option = "DENY"
      override     = false
    }
    referrer_policy {
      referrer_policy = "strict-origin-when-cross-origin"
      override        = false
    }
    strict_transport_security {
      access_control_max_age_sec = 63072000
      include_subdomains         = true
      override                   = false
    }
  }
}

resource "aws_s3_bucket_cors_configuration" "website" {
  count  = var.enable_cors ? 1 : 0
  bucket = aws_s3_bucket.website.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

#########################
### DNS configuration ###
#########################
resource "aws_route53_record" "website" {
  for_each = toset(local.full_addresses)

  zone_id = var.route53_zone_id
  name    = each.value
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.website.domain_name
    zone_id                = aws_cloudfront_distribution.website.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "website_ipv6" {
  for_each = toset(local.full_addresses)

  zone_id = var.route53_zone_id
  name    = each.value
  type    = "AAAA"

  alias {
    name                   = aws_cloudfront_distribution.website.domain_name
    zone_id                = aws_cloudfront_distribution.website.hosted_zone_id
    evaluate_target_health = false
  }
}

data "aws_route53_zone" "selected" {
  zone_id = var.route53_zone_id
}
