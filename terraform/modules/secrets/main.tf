resource "aws_secretsmanager_secret" "database_url" {
  name = "${var.app_name}/${var.environment}/database-url"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "database_url" {
  secret_id     = aws_secretsmanager_secret.database_url.id
  secret_string = "postgresql://chatadmin:${var.db_password}@PLACEHOLDER:5432/chat_db?sslmode=require"
}

resource "aws_secretsmanager_secret" "jwt_secret" {
  name = "${var.app_name}/${var.environment}/jwt-secret"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "jwt_secret" {
  secret_id     = aws_secretsmanager_secret.jwt_secret.id
  secret_string = var.jwt_secret
}

resource "aws_secretsmanager_secret" "stream_api_key" {
  name = "${var.app_name}/${var.environment}/stream-api-key"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "stream_api_key" {
  secret_id     = aws_secretsmanager_secret.stream_api_key.id
  secret_string = var.stream_api_key
}

resource "aws_secretsmanager_secret" "stream_api_secret" {
  name = "${var.app_name}/${var.environment}/stream-api-secret"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "stream_api_secret" {
  secret_id     = aws_secretsmanager_secret.stream_api_secret.id
  secret_string = var.stream_api_secret
}
