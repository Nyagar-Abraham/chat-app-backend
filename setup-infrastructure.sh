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

# Check if security features are enabled
HTTPS_ENABLED=$(terraform output -raw enable_https 2>/dev/null || echo "false")
WAF_ENABLED=$(terraform output -raw enable_waf 2>/dev/null || echo "false")
GUARDDUTY_ENABLED=$(terraform output -raw enable_guardduty 2>/dev/null || echo "false")

echo ""
echo "✅ Infrastructure setup complete!"
echo ""
echo "📊 Infrastructure Summary:"
echo "  ECR Repository: $ECR_REPO"
echo "  Database Endpoint: $DB_ENDPOINT"
echo "  Load Balancer: $ALB_DNS"

if [[ "$HTTPS_ENABLED" = "true" ]]; then
    CERT_ARN=$(terraform output -raw certificate_arn 2>/dev/null)
    echo "  SSL Certificate: $CERT_ARN"
    echo "  Domain: hxrrvpsxtz.xyz"
    echo ""
    echo "⚠️  HTTPS ENABLED - DNS Validation Required:"
    echo "  1. Add DNS validation records to GoDaddy"
    echo "  2. View records: terraform output certificate_validation_records"
    echo "  3. Wait for validation (5-30 minutes)"
    echo "  4. Add A/CNAME records pointing to ALB"
fi

if [[ "$WAF_ENABLED" = "true" ]]; then
    WAF_ID=$(terraform output -raw waf_web_acl_id 2>/dev/null)
    echo "  WAF Web ACL: $WAF_ID"
fi

if [[ "$GUARDDUTY_ENABLED" = "true" ]]; then
    GUARDDUTY_ID=$(terraform output -raw guardduty_detector_id 2>/dev/null)
    echo "  GuardDuty Detector: $GUARDDUTY_ID"
fi

echo ""
echo "🔄 Next Steps:"
echo "  1. Update database URL secret:"
echo "     aws secretsmanager update-secret --secret-id chat-app/prod/database-url --secret-string 'postgresql://chatadmin:PASSWORD@${DB_ENDPOINT}/chat_db?sslmode=require'"
if [ "$HTTPS_ENABLED" = "true" ]; then
    echo "  2. Add DNS validation records to GoDaddy (see SECURITY_DEPLOYMENT_GUIDE.md)"
    echo "  3. Confirm security alert email subscription"
    echo "  4. Update GitHub Secrets with AWS credentials"
    echo "  5. Push code to trigger CI/CD deployment"
else
    echo "  2. Update GitHub Secrets with AWS credentials"
    echo "  3. Push code to trigger CI/CD deployment"
fi
echo ""
echo "📚 Documentation:"
echo "  - Security setup: terraform/SECURITY_DEPLOYMENT_GUIDE.md"
echo "  - Security commands: terraform/SECURITY_COMMANDS.md"
echo "  - IAM user setup: terraform/IAM_USER_SETUP.md"