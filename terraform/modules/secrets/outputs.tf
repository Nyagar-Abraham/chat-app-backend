output "database_url_secret_arn" {
  value = aws_secretsmanager_secret.database_url.arn
}

output "jwt_secret_arn" {
  value = aws_secretsmanager_secret.jwt_secret.arn
}

output "stream_api_key_secret_arn" {
  value = aws_secretsmanager_secret.stream_api_key.arn
}

output "stream_api_secret_secret_arn" {
  value = aws_secretsmanager_secret.stream_api_secret.arn
}
