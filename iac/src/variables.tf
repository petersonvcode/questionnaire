variable "environment" {
  type        = string
  description = "The environment of the website"
  validation {
    condition     = contains(["dev", "prod"], var.environment)
    error_message = "Invalid environment. Valid values are: dev, prod."
  }
}

variable "route53_zone_id" {
  type        = string
  description = "The Route53 hosted zone ID for the root domain"
}

variable "cost_mode" {
  type        = string
  description = "The cost mode of the website"
  default     = "standard"
  validation {
    condition     = contains(["standard", "cheap"], var.cost_mode)
    error_message = "Invalid cost mode. Valid values are: standard, cheap."
  }
}

variable "acm_certificate_arn" {
  type        = string
  description = "The ACM certificate ARN for the website"
}

variable "vps_ip" {
  type        = string
  description = "The IP address of the VPS"
}