# SonarCloud Setup Guide

Complete guide to integrate SonarCloud code quality analysis into your chat-app CI/CD pipeline.

## 🎯 What is SonarCloud?

SonarCloud provides:
- **Code Quality Analysis** - Detect bugs, vulnerabilities, and code smells
- **Security Scanning** - Identify security hotspots
- **Code Coverage** - Track test coverage metrics
- **Technical Debt** - Measure maintainability
- **Quality Gates** - Enforce quality standards

## 📋 Prerequisites

- GitHub repository (public or private)
- GitHub account with admin access

## 🚀 Setup Steps

### Step 1: Sign Up for SonarCloud

1. Go to [SonarCloud](https://sonarcloud.io/)
2. Click **"Log in"** → **"With GitHub"**
3. Authorize SonarCloud to access your GitHub account

### Step 2: Import Your Repository

1. Click **"+"** → **"Analyze new project"**
2. Select your organization or create one:
   - Organization key: `nyagar-abraham` (or your GitHub username)
3. Select repository: `chat-app`
4. Click **"Set Up"**

### Step 3: Configure Analysis Method

1. Choose **"With GitHub Actions"**
2. SonarCloud will show you the configuration

### Step 4: Get Your SonarCloud Token

1. Go to **My Account** → **Security**
2. Generate a new token:
   - Name: `chat-app-github-actions`
   - Type: **User Token**
   - Expiration: **No expiration** (or 90 days)
3. **Copy the token** (you won't see it again!)

### Step 5: Add GitHub Secret

1. Go to your GitHub repository
2. **Settings** → **Secrets and variables** → **Actions**
3. Click **"New repository secret"**
4. Add secret:
   - Name: `SONAR_TOKEN`
   - Value: `<paste-your-sonarcloud-token>`
5. Click **"Add secret"**

### Step 6: Update sonar-project.properties

Edit `sonar-project.properties` and update:

```properties
# Replace with your actual values
sonar.projectKey=YOUR_GITHUB_USERNAME_chat-app
sonar.organization=YOUR_GITHUB_USERNAME
```

To find these values:
1. Go to SonarCloud → Your Project
2. **Project Information** shows both values

### Step 7: Commit and Push

```bash
git add .github/workflows/deploy.yml sonar-project.properties SONARCLOUD_SETUP.md
git commit -m "Add SonarCloud code quality analysis"
git push origin main
```

### Step 8: Verify Integration

1. Go to **Actions** tab in GitHub
2. Watch the workflow run
3. Check the **SonarCloud Scan** step
4. Go to [SonarCloud Dashboard](https://sonarcloud.io/projects) to see results

## 📊 Understanding SonarCloud Results

### Quality Gate

Your code must pass these checks:
- **Coverage**: > 80% (configurable)
- **Duplications**: < 3%
- **Maintainability Rating**: A
- **Reliability Rating**: A
- **Security Rating**: A

### Metrics Explained

- **Bugs**: Logic errors that will cause issues
- **Vulnerabilities**: Security weaknesses
- **Code Smells**: Maintainability issues
- **Coverage**: % of code tested
- **Duplications**: Repeated code blocks
- **Technical Debt**: Time to fix all issues

## 🔧 Configuration Options

### Adjust Quality Gate

In `sonar-project.properties`:

```properties
# Require 80% coverage
sonar.coverage.minimum=80

# Allow up to 5% duplication
sonar.cpd.exclusions=**/*_test.go
```

### Exclude Files

```properties
sonar.exclusions=**/vendor/**,**/testdata/**,terraform/**,**/*.pb.go
```

### Add Custom Rules

1. Go to SonarCloud → **Quality Profiles**
2. Create custom profile for Go
3. Activate/deactivate specific rules

## 🎨 Add SonarCloud Badge to README

Add this to your `README.md`:

```markdown
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=YOUR_PROJECT_KEY&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=YOUR_PROJECT_KEY)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=YOUR_PROJECT_KEY&metric=coverage)](https://sonarcloud.io/summary/new_code?id=YOUR_PROJECT_KEY)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=YOUR_PROJECT_KEY&metric=bugs)](https://sonarcloud.io/summary/new_code?id=YOUR_PROJECT_KEY)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=YOUR_PROJECT_KEY&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=YOUR_PROJECT_KEY)
```

Replace `YOUR_PROJECT_KEY` with your actual project key (e.g., `Nyagar-Abraham_chat-app`).

## 🚨 Troubleshooting

### "Quality Gate Failed"

Check which metric failed:
1. Go to SonarCloud dashboard
2. Click on your project
3. Review **New Code** tab
4. Fix issues and push again

### "Coverage Too Low"

Add more tests:
```bash
# Check current coverage
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### "Authentication Error"

1. Verify `SONAR_TOKEN` secret is set correctly
2. Check token hasn't expired
3. Regenerate token if needed

### "Project Not Found"

Update `sonar-project.properties`:
```properties
sonar.projectKey=YOUR_ACTUAL_PROJECT_KEY
sonar.organization=YOUR_ACTUAL_ORG
```

## 🔒 Security Best Practices

1. **Never commit SONAR_TOKEN** to repository
2. Use **GitHub Secrets** for sensitive data
3. Set token expiration (90 days recommended)
4. Rotate tokens regularly
5. Use **branch protection** to enforce quality gates

## 📈 Advanced Features

### Pull Request Analysis

SonarCloud automatically analyzes PRs:
- Shows new issues introduced
- Comments on PR with findings
- Blocks merge if quality gate fails

### Branch Analysis

Track quality across branches:
```yaml
# In deploy.yml, add for feature branches
on:
  pull_request:
    branches: [main]
```

### Custom Metrics

Track custom metrics:
```properties
sonar.custom.metric.name=custom_value
```

## 💰 Pricing

- **Public repositories**: FREE forever
- **Private repositories**: 
  - Free for open source
  - Paid plans start at $10/month

## 📚 Additional Resources

- [SonarCloud Documentation](https://docs.sonarcloud.io/)
- [Go Language Support](https://docs.sonarcloud.io/enriching/languages/go/)
- [Quality Gates](https://docs.sonarcloud.io/improving/quality-gates/)
- [GitHub Integration](https://docs.sonarcloud.io/getting-started/github/)

## ✅ Verification Checklist

- [ ] SonarCloud account created
- [ ] Repository imported to SonarCloud
- [ ] SONAR_TOKEN added to GitHub Secrets
- [ ] sonar-project.properties configured
- [ ] Workflow updated and pushed
- [ ] First analysis completed successfully
- [ ] Quality gate passing
- [ ] Badges added to README (optional)

## 🎯 Next Steps

1. **Review initial scan results** on SonarCloud
2. **Fix critical issues** (bugs, vulnerabilities)
3. **Improve test coverage** to > 80%
4. **Set up branch protection** to require quality gate
5. **Monitor trends** over time

---

**Need help?** Check the [SonarCloud Community](https://community.sonarsource.com/) or open an issue.
