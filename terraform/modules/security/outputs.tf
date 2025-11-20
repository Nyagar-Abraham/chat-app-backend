output "vpc_flow_logs_log_group" {
  description = "CloudWatch log group for VPC flow logs"
  value       = var.enable_vpc_flow_logs ? aws_cloudwatch_log_group.vpc_flow_logs[0].name : null
}

output "guardduty_detector_id" {
  description = "ID of the GuardDuty detector"
  value       = var.enable_guardduty ? aws_guardduty_detector.main[0].id : null
}

output "security_alerts_topic_arn" {
  description = "ARN of the SNS topic for security alerts"
  value       = var.enable_security_alerts ? aws_sns_topic.security_alerts[0].arn : null
}

output "config_recorder_id" {
  description = "ID of the AWS Config recorder"
  value       = var.enable_config ? aws_config_configuration_recorder.main[0].id : null
}

output "config_bucket_name" {
  description = "Name of the S3 bucket for AWS Config"
  value       = var.enable_config ? aws_s3_bucket.config[0].id : null
}
