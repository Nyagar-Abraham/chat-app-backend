variable "app_name" {
  type = string
}

variable "environment" {
  type = string
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "stream_api_key" {
  type      = string
  sensitive = true
}

variable "stream_api_secret" {
  type      = string
  sensitive = true
}
