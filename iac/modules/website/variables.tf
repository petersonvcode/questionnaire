variable "name" {
  type        = string
  description = "The name of the website"
}

variable "subdomains" {
  type        = list(string)
  description = "subdomains of the website"
  default     = []
}

variable "cost_mode" {
  type        = string
  description = "The cost mode of the website"
  default     = "standard"
}

variable "acm_certificate_arn" {
  type        = string
  description = "The ACM certificate ARN for the website"
}

variable "environment" {
  type        = string
  description = "The environment of the website"
}

variable "tags" {
  type        = map(string)
  description = "The tags of the website resources"
  default     = {}
}

variable "enable_cors" {
  type        = bool
  description = "Whether to enable CORS"
  default     = false
}

variable "route53_zone_id" {
  type        = string
  description = "The Route53 hosted zone ID for the root domain"
}