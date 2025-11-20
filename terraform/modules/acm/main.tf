# ACM Certificate for HTTPS
resource "aws_acm_certificate" "main" {
  domain_name       = var.domain_name
  validation_method = "DNS"

  subject_alternative_names = var.subject_alternative_names

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name        = "${var.app_name}-certificate"
    Environment = var.environment
  }
}

# Note: Since you're using GoDaddy, you'll need to manually add the DNS validation records
# The records needed will be output after terraform apply
