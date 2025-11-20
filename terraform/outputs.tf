output "alb_dns_name" {
  description = "ALB DNS name"
  value       = module.alb.alb_dns_name
}

output "db_endpoint" {
  description = "RDS endpoint"
  value       = module.rds.db_endpoint
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = module.ecs.cluster_name
}

output "ecr_repository_url" {
  description = "ECR repository URL"
  value       = module.ecs.ecr_repository_url
}

# Security Outputs
output "certificate_arn" {
  description = "ACM certificate ARN (if HTTPS enabled)"
  value       = var.enable_https ? module.acm[0].certificate_arn : null
}

output "certificate_validation_records" {
  description = "DNS validation records for GoDaddy (if HTTPS enabled)"
  value       = var.enable_https ? module.acm[0].domain_validation_options : null
}

output "waf_web_acl_id" {
  description = "WAF Web ACL ID (if WAF enabled)"
  value       = var.enable_waf ? module.waf[0].web_acl_id : null
}

output "waf_log_group" {
  description = "WAF CloudWatch log group (if WAF enabled)"
  value       = var.enable_waf ? module.waf[0].log_group_name : null
}

output "security_alert_topic" {
  description = "SNS topic for security alerts (if enabled)"
  value       = var.enable_security_alerts ? module.security.security_alerts_topic_arn : null
}

output "guardduty_detector_id" {
  description = "GuardDuty detector ID (if enabled)"
  value       = var.enable_guardduty ? module.security.guardduty_detector_id : null
}
