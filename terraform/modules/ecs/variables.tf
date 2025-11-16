variable "app_name" {
  type = string
}

variable "environment" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "alb_target_group_arn" {
  type = string
}

variable "alb_security_group_id" {
  type = string
}

variable "database_url_secret_arn" {
  type = string
}

variable "jwt_secret_arn" {
  type = string
}

variable "stream_api_key_secret_arn" {
  type = string
}

variable "stream_api_secret_secret_arn" {
  type = string
}

variable "db_endpoint" {
  type = string
}
