variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "app_name" {
  description = "Application name"
  type        = string
  default     = "chat-app"
}

variable "environment" {
  description = "Environment (dev/prod)"
  type        = string
  default     = "prod"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "JWT secret key"
  type        = string
  sensitive   = true
}

variable "stream_api_key" {
  description = "Stream Chat API key"
  type        = string
  sensitive   = true
}

variable "stream_api_secret" {
  description = "Stream Chat API secret"
  type        = string
  sensitive   = true
}

variable "docker_image" {
  description = "Docker image URI"
  type        = string
}

# HTTPS/SSL Configuration
variable "enable_https" {
  description = "Enable HTTPS with SSL certificate"
  type        = bool
  default     = false
}

variable "domain_name" {
  description = "Primary domain name for SSL certificate"
  type        = string
  default     = ""
}

variable "subject_alternative_names" {
  description = "Additional domain names (e.g., www subdomain)"
  type        = list(string)
  default     = []
}

# WAF Configuration
variable "enable_waf" {
  description = "Enable Web Application Firewall"
  type        = bool
  default     = false
}

variable "waf_rate_limit" {
  description = "WAF rate limit per IP (requests per 5 minutes)"
  type        = number
  default     = 2000
}

variable "waf_enable_geo_blocking" {
  description = "Enable geographic blocking in WAF"
  type        = bool
  default     = false
}

variable "waf_blocked_countries" {
  description = "List of country codes to block (e.g., ['CN', 'RU'])"
  type        = list(string)
  default     = []
}

variable "waf_log_retention_days" {
  description = "WAF log retention in days"
  type        = number
  default     = 30
}

# VPC Flow Logs Configuration
variable "enable_vpc_flow_logs" {
  description = "Enable VPC Flow Logs"
  type        = bool
  default     = true
}

variable "flow_logs_retention_days" {
  description = "VPC Flow Logs retention in days"
  type        = number
  default     = 30
}

variable "flow_logs_traffic_type" {
  description = "Traffic type to log (ACCEPT, REJECT, ALL)"
  type        = string
  default     = "ALL"
}

# GuardDuty Configuration
variable "enable_guardduty" {
  description = "Enable AWS GuardDuty threat detection"
  type        = bool
  default     = true
}

variable "guardduty_enable_s3_logs" {
  description = "Enable S3 protection in GuardDuty"
  type        = bool
  default     = true
}

variable "guardduty_enable_malware_protection" {
  description = "Enable malware protection in GuardDuty"
  type        = bool
  default     = true
}

variable "guardduty_finding_frequency" {
  description = "GuardDuty finding frequency"
  type        = string
  default     = "FIFTEEN_MINUTES"
}

# Security Alerts Configuration
variable "enable_security_alerts" {
  description = "Enable security alerts via SNS"
  type        = bool
  default     = true
}

variable "security_alert_email" {
  description = "Email for security alerts"
  type        = string
  default     = ""
}

variable "guardduty_alert_severities" {
  description = "GuardDuty severity levels to alert on (4.0+ = Medium, High, Critical)"
  type        = list(number)
  default     = [4, 4.0, 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 4.8, 4.9, 5, 5.0, 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7, 5.8, 5.9, 6, 6.0, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8, 6.9, 7, 7.0, 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8, 7.9, 8, 8.0, 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7, 8.8, 8.9]
}

# AWS Config
variable "enable_config" {
  description = "Enable AWS Config for compliance monitoring"
  type        = bool
  default     = false
}
