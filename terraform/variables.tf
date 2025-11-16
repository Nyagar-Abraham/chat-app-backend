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
