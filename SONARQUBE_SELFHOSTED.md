# Self-Hosted SonarQube Setup (Alternative)

Guide for running your own SonarQube server instead of using SonarCloud.

## 🎯 When to Use Self-Hosted

- Private enterprise projects
- Need full control over data
- Custom plugins required
- Air-gapped environments

## 🚀 Quick Setup with Docker

### Option 1: Docker Compose (Recommended)

Create `sonarqube-docker-compose.yml`:

```yaml
version: "3"

services:
  sonarqube:
    image: sonarqube:community
    depends_on:
      - db
    environment:
      SONAR_JDBC_URL: jdbc:postgresql://db:5432/sonar
      SONAR_JDBC_USERNAME: sonar
      SONAR_JDBC_PASSWORD: sonar
    volumes:
      - sonarqube_data:/opt/sonarqube/data
      - sonarqube_extensions:/opt/sonarqube/extensions
      - sonarqube_logs:/opt/sonarqube/logs
    ports:
      - "9000:9000"

  db:
    image: postgres:15
    environment:
      POSTGRES_USER: sonar
      POSTGRES_PASSWORD: sonar
      POSTGRES_DB: sonar
    volumes:
      - postgresql:/var/lib/postgresql
      - postgresql_data:/var/lib/postgresql/data

volumes:
  sonarqube_data:
  sonarqube_extensions:
  sonarqube_logs:
  postgresql:
  postgresql_data:
```

Start SonarQube:

```bash
docker-compose -f sonarqube-docker-compose.yml up -d
```

Access at: http://localhost:9000
- Default login: `admin` / `admin`

### Option 2: AWS EC2 Deployment

```bash
# Launch EC2 instance (t3.medium minimum)
# Install Docker
sudo yum update -y
sudo yum install docker -y
sudo service docker start

# Run SonarQube
docker run -d --name sonarqube \
  -p 9000:9000 \
  -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true \
  sonarqube:community
```

## 🔧 GitHub Actions Configuration

Update `.github/workflows/deploy.yml`:

```yaml
- name: SonarQube Scan
  uses: sonarsource/sonarqube-scan-action@master
  env:
    SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
    SONAR_HOST_URL: ${{ secrets.SONAR_HOST_URL }}  # Your server URL
```

Update `sonar-project.properties`:

```properties
# Remove SonarCloud specific settings
# sonar.organization=nyagar-abraham  # Remove this line

# Keep these
sonar.projectKey=chat-app
sonar.projectName=Chat App Backend
sonar.sources=.
sonar.go.coverage.reportPaths=coverage.out
```

## 🔐 Setup Steps

1. **Access SonarQube**: http://your-server:9000
2. **Login**: admin/admin (change password)
3. **Create Project**:
   - Name: `chat-app`
   - Key: `chat-app`
4. **Generate Token**:
   - My Account → Security → Generate Token
   - Name: `github-actions`
5. **Add GitHub Secrets**:
   - `SONAR_TOKEN`: Your generated token
   - `SONAR_HOST_URL`: http://your-server:9000

## 💰 Cost Comparison

**SonarCloud (Managed)**:
- Free for public repos
- $10/month for private repos

**Self-Hosted**:
- EC2 t3.medium: ~$30/month
- Storage: ~$5/month
- Total: ~$35/month + maintenance time

## ✅ Recommendation

Use **SonarCloud** unless you have specific requirements for self-hosting.

---

For most projects, follow `SONARCLOUD_SETUP.md` instead.
