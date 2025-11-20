# IAM User Setup for Terraform Deployment

This guide helps you create an IAM user with the necessary permissions to deploy the entire chat-app infrastructure using Terraform.

## 🎯 Overview

The IAM user will have permissions to manage:
- VPC, Subnets, Internet Gateways, NAT Gateways
- RDS PostgreSQL databases
- ECS Fargate clusters and services
- ECR repositories
- Application Load Balancers
- ACM SSL certificates
- WAF Web ACLs
- Secrets Manager
- GuardDuty
- VPC Flow Logs
- CloudWatch Logs
- SNS topics
- S3 buckets (for Terraform state and AWS Config)

## 📋 Prerequisites

- AWS account with admin access (to create the IAM user)
- AWS CLI installed and configured

## 🚀 Option 1: Using AWS Console (Recommended for Beginners)

### Step 1: Create IAM Policies (3 policies)

AWS has a 6144 character limit per policy, so we split into 3 policies:

**Policy 1: Compute Resources**
1. Log in to [AWS IAM Console](https://console.aws.amazon.com/iam/)
2. Click **Policies** → **Create policy** → **JSON** tab
3. Copy contents of `iam-policy-1-compute.json` and paste
4. Name: `ChatAppTerraformCompute`
5. Click **Create policy**

**Policy 2: Security Services**
1. Click **Create policy** → **JSON** tab
2. Copy contents of `iam-policy-2-security.json` and paste
3. Name: `ChatAppTerraformSecurity`
4. Click **Create policy**

**Policy 3: Monitoring & Storage**
1. Click **Create policy** → **JSON** tab
2. Copy contents of `iam-policy-3-monitoring.json` and paste
3. Name: `ChatAppTerraformMonitoring`
4. Click **Create policy**

### Step 2: Create IAM User

1. In IAM Console, click **Users** → **Add users**
2. Username: `terraform-chat-app`
3. Select **Access key - Programmatic access**
4. Click **Next: Permissions**
5. Select **Attach existing policies directly**
6. Search and check ALL three policies:
   - `ChatAppTerraformCompute`
   - `ChatAppTerraformSecurity`
   - `ChatAppTerraformMonitoring`
7. Click **Next: Tags** → **Next: Review** → **Create user**

### Step 3: Save Credentials

⚠️ **IMPORTANT**: Save these credentials immediately - you won't see them again!

1. Click **Download .csv** to save credentials
2. Or copy:
   - Access key ID
   - Secret access key

### Step 4: Configure AWS CLI

```bash
aws configure --profile terraform-chat-app
```

Enter:
- AWS Access Key ID: `<your-access-key-id>`
- AWS Secret Access Key: `<your-secret-access-key>`
- Default region: `us-east-1`
- Default output format: `json`

### Step 5: Test Access

```bash
# Test with the new profile
aws sts get-caller-identity --profile terraform-chat-app

# Set as default for current session
export AWS_PROFILE=terraform-chat-app
```

## 🚀 Option 2: Using AWS CLI (Faster)

### Step 1: Create IAM Policies

```bash
# Policy 1: Compute
aws iam create-policy \
  --policy-name ChatAppTerraformCompute \
  --policy-document file://iam-policy-1-compute.json

# Policy 2: Security
aws iam create-policy \
  --policy-name ChatAppTerraformSecurity \
  --policy-document file://iam-policy-2-security.json

# Policy 3: Monitoring
aws iam create-policy \
  --policy-name ChatAppTerraformMonitoring \
  --policy-document file://iam-policy-3-monitoring.json
```

### Step 2: Create IAM User

```bash
aws iam create-user --user-name terraform-chat-app
```

### Step 3: Attach Policies to User

```bash
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Attach all 3 policies
aws iam attach-user-policy \
  --user-name terraform-chat-app \
  --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformCompute

aws iam attach-user-policy \
  --user-name terraform-chat-app \
  --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformSecurity

aws iam attach-user-policy \
  --user-name terraform-chat-app \
  --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformMonitoring
```

### Step 4: Create Access Keys

```bash
aws iam create-access-key --user-name terraform-chat-app
```

⚠️ **Save the output immediately!**

### Step 5: Configure AWS CLI

```bash
# Add to AWS credentials
aws configure --profile terraform-chat-app

# Or manually add to ~/.aws/credentials
cat >> ~/.aws/credentials << EOF

[terraform-chat-app]
aws_access_key_id = YOUR_ACCESS_KEY_ID
aws_secret_access_key = YOUR_SECRET_ACCESS_KEY
EOF

# Add to ~/.aws/config
cat >> ~/.aws/config << EOF

[profile terraform-chat-app]
region = us-east-1
output = json
EOF
```

## 🔒 Security Best Practices

### 1. Use MFA (Multi-Factor Authentication)

Enable MFA for the IAM user:

```bash
# In AWS Console
1. Go to IAM → Users → terraform-chat-app
2. Click "Security credentials" tab
3. Click "Manage" next to "Assigned MFA device"
4. Follow the setup wizard
```

### 2. Rotate Access Keys Regularly

```bash
# Create new access key
aws iam create-access-key --user-name terraform-chat-app

# Update your credentials
aws configure --profile terraform-chat-app

# Delete old access key
aws iam delete-access-key \
  --user-name terraform-chat-app \
  --access-key-id OLD_ACCESS_KEY_ID
```

### 3. Use Environment Variables (Alternative)

Instead of AWS profiles, use environment variables:

```bash
export AWS_ACCESS_KEY_ID="your-access-key-id"
export AWS_SECRET_ACCESS_KEY="your-secret-access-key"
export AWS_DEFAULT_REGION="us-east-1"
```

### 4. Store Credentials Securely

- Never commit credentials to Git
- Use a password manager for storing credentials
- Consider using AWS Secrets Manager or Parameter Store

## 🧪 Verify Permissions

Test that the user has the required permissions:

```bash
# Set profile
export AWS_PROFILE=terraform-chat-app

# Test VPC permissions
aws ec2 describe-vpcs --region us-east-1

# Test RDS permissions
aws rds describe-db-instances --region us-east-1

# Test ECS permissions
aws ecs list-clusters --region us-east-1

# Test Secrets Manager permissions
aws secretsmanager list-secrets --region us-east-1

# Test S3 permissions (Terraform state)
aws s3 ls s3://chat-app-terraform-state-869935099753
```

## 🚀 Deploy Infrastructure

Now you can deploy the infrastructure:

```bash
cd terraform

# Set the profile
export AWS_PROFILE=terraform-chat-app

# Initialize Terraform
terraform init

# Plan
terraform plan

# Apply
terraform apply
```

## 🔧 Troubleshooting

### "Access Denied" Errors

1. Verify the policy is attached:
```bash
aws iam list-attached-user-policies --user-name terraform-chat-app
```

2. Check if you're using the correct profile:
```bash
aws sts get-caller-identity
```

3. Verify the policy has all required permissions:
```bash
aws iam get-policy-version \
  --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformPolicy \
  --version-id v1
```

### "Invalid Credentials" Errors

1. Verify credentials are correct:
```bash
cat ~/.aws/credentials | grep -A 3 terraform-chat-app
```

2. Test credentials:
```bash
aws sts get-caller-identity --profile terraform-chat-app
```

### Missing Permissions

If you encounter missing permissions:

1. Identify which policy needs the permission
2. Update the appropriate JSON file
3. Create new policy version:
```bash
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Update the relevant policy (example for compute)
aws iam create-policy-version \
  --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformCompute \
  --policy-document file://iam-policy-1-compute.json \
  --set-as-default
```

## 🗑️ Cleanup (When No Longer Needed)

### Delete IAM User

```bash
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Delete access keys
aws iam list-access-keys --user-name terraform-chat-app
aws iam delete-access-key \
  --user-name terraform-chat-app \
  --access-key-id ACCESS_KEY_ID

# Detach all policies
aws iam detach-user-policy \
  --user-name terraform-chat-app \
  --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformCompute

aws iam detach-user-policy \
  --user-name terraform-chat-app \
  --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformSecurity

aws iam detach-user-policy \
  --user-name terraform-chat-app \
  --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformMonitoring

# Delete user
aws iam delete-user --user-name terraform-chat-app
```

### Delete IAM Policies (Optional)

```bash
aws iam delete-policy --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformCompute
aws iam delete-policy --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformSecurity
aws iam delete-policy --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/ChatAppTerraformMonitoring
```

## 📚 Additional Resources

- [AWS IAM Best Practices](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html)
- [Terraform AWS Provider Authentication](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#authentication-and-configuration)
- [AWS CLI Configuration](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)

## ✅ Checklist

- [ ] IAM policy created
- [ ] IAM user created
- [ ] Policy attached to user
- [ ] Access keys generated and saved securely
- [ ] AWS CLI configured with new credentials
- [ ] Permissions verified with test commands
- [ ] MFA enabled (recommended)
- [ ] Credentials stored securely (not in Git)

---

**Ready to deploy?** Run `./deploy-security.sh` or `terraform apply` to start!
