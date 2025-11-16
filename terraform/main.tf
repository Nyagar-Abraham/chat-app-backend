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

# ALB Module
module "alb" {
  source = "./modules/alb"

  app_name          = var.app_name
  environment       = var.environment
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
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
