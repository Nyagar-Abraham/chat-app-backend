#!/bin/bash
set -e

echo "🏗️ Infrastructure Setup Script"
echo "=============================="
echo "This script sets up AWS infrastructure only."
echo "Actual deployment happens via GitHub Actions CI/CD."
echo ""

# Check prerequisites
command -v aws >/dev/null 2>&1 || { echo "❌ AWS CLI not installed"; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo "❌ Terraform not installed"; exit 1; }

echo "✅ Prerequisites check passed"

# Step 1: Create S3 bucket for Terraform state
echo ""
echo "📦 Step 1: Creating S3 bucket for Terraform state..."
aws s3api create-bucket --bucket chat-app-terraform-state-869935099753 --region us-east-1 2>/dev/null || echo "Bucket already exists"
aws s3api put-bucket-versioning --bucket chat-app-terraform-state-869935099753 --versioning-configuration Status=Enabled

# Step 2: Initialize Terraform
echo ""
echo "🔧 Step 2: Initializing Terraform..."
cd terraform
terraform init

# Step 3: Plan infrastructure
echo ""
echo "📋 Step 3: Planning infrastructure..."
terraform plan

# Step 4: Apply infrastructure
echo ""
echo "🏗️ Step 4: Deploying infrastructure..."
echo "This will create VPC, RDS, ECS, ALB, and other resources"
read -p "Continue? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "Infrastructure setup cancelled"
    exit 0
fi

terraform apply

# Step 5: Save outputs
echo ""
echo "📝 Step 5: Saving infrastructure outputs..."
terraform output > ../infrastructure-outputs.txt
ECR_REPO=$(terraform output -raw ecr_repository_url)
DB_ENDPOINT=$(terraform output -raw db_endpoint)
ALB_DNS=$(terraform output -raw alb_dns_name)

echo ""
echo "✅ Infrastructure setup complete!"
echo ""
echo "📊 Infrastructure Summary:"
echo "  ECR Repository: $ECR_REPO"
echo "  Database Endpoint: $DB_ENDPOINT"
echo "  Load Balancer: $ALB_DNS"
echo ""
echo "🔄 Next Steps:"
echo "  1. Update GitHub Secrets with AWS credentials"
echo "  2. Update database URL secret with actual endpoint"
echo "  3. Push code to trigger CI/CD deployment"
echo ""
echo "📝 Manual secret update command:"
echo "  aws secretsmanager update-secret --secret-id chat-app/prod/database-url --secret-string 'postgresql://chatadmin:PASSWORD@${DB_ENDPOINT}/chat_db?sslmode=require'"