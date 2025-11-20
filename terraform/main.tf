terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  backend "s3" {
    bucket = "chat-app-terraform-state-869935099753"
    key    = "prod/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC Module
module "vpc" {
  source = "./modules/vpc"
  
  app_name    = var.app_name
  environment = var.environment
}

# Secrets Module
module "secrets" {

  source = "./modules/secrets"

  app_name           = var.app_name
  environment        = var.environment
  db_password        = var.db_password
  jwt_secret         = var.jwt_secret
  stream_api_key     = var.stream_api_key
  stream_api_secret  = var.stream_api_secret
}

# RDS Module
module "rds" {
  source = "./modules/rds"

  app_name           = var.app_name
  environment        = var.environment
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  db_password        = var.db_password
}

# ACM Module (SSL/TLS Certificate)
module "acm" {
  count  = var.enable_https ? 1 : 0
  source = "./modules/acm"

  app_name                  = var.app_name
  environment               = var.environment
  domain_name               = var.domain_name
  subject_alternative_names = var.subject_alternative_names
}

# WAF Module (Web Application Firewall)
module "waf" {
  count  = var.enable_waf ? 1 : 0
  source = "./modules/waf"

  app_name            = var.app_name
  environment         = var.environment
  rate_limit          = var.waf_rate_limit
  enable_geo_blocking = var.waf_enable_geo_blocking
  blocked_countries   = var.waf_blocked_countries
  log_retention_days  = var.waf_log_retention_days
}

# ALB Module
module "alb" {
  source = "./modules/alb"

  app_name          = var.app_name
  environment       = var.environment
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids

  # HTTPS Configuration
  enable_https    = var.enable_https
  certificate_arn = var.enable_https ? module.acm[0].certificate_arn : ""

  # WAF Configuration
  enable_waf  = var.enable_waf
  waf_acl_arn = var.enable_waf ? module.waf[0].web_acl_arn : ""

  depends_on = [
    module.acm,
    module.waf
  ]
}

# ECS Module
module "ecs" {
  source = "./modules/ecs"

  app_name              = var.app_name
  environment           = var.environment
  vpc_id                = module.vpc.vpc_id
  private_subnet_ids    = module.vpc.private_subnet_ids
  alb_target_group_arn  = module.alb.target_group_arn
  alb_security_group_id = module.alb.alb_security_group_id

  database_url_secret_arn       = module.secrets.database_url_secret_arn
  jwt_secret_arn                = module.secrets.jwt_secret_arn
  stream_api_key_secret_arn     = module.secrets.stream_api_key_secret_arn
  stream_api_secret_secret_arn  = module.secrets.stream_api_secret_secret_arn

  db_endpoint = module.rds.db_endpoint
}

# Security Module (VPC Flow Logs, GuardDuty, etc.)
module "security" {
  source = "./modules/security"

  app_name    = var.app_name
  environment = var.environment
  vpc_id      = module.vpc.vpc_id

  # VPC Flow Logs Configuration
  enable_vpc_flow_logs      = var.enable_vpc_flow_logs
  flow_logs_retention_days  = var.flow_logs_retention_days
  flow_logs_traffic_type    = var.flow_logs_traffic_type

  # GuardDuty Configuration
  enable_guardduty                     = var.enable_guardduty
  guardduty_enable_s3_logs             = var.guardduty_enable_s3_logs
  guardduty_enable_malware_protection  = var.guardduty_enable_malware_protection
  guardduty_finding_frequency          = var.guardduty_finding_frequency

  # Security Alerts Configuration
  enable_security_alerts    = var.enable_security_alerts
  alert_email               = var.security_alert_email
  guardduty_alert_severities = var.guardduty_alert_severities

  # AWS Config
  enable_config = var.enable_config
}
