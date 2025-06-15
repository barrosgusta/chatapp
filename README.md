# ChatApp

A full-stack chat application with a React + Vite + TypeScript frontend and Go-based backend microservices. The backend uses AWS DynamoDB and SQS for message storage and delivery. This project supports both local development (via Docker Compose) and cloud deployment (via Helm and Terraform).

---

## Table of Contents

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Local Development](#local-development)
- [Environment Variables](#environment-variables)
- [Cloud Deployment](#cloud-deployment)
- [Step-by-Step Cloud Deployment Guide](#step-by-step-cloud-deployment-guide)
- [Project Structure](#project-structure)

---

## Architecture

- **Frontend:** React + Vite + TypeScript (`chatapp-frontend`)
- **Backend:**
  - `gateway-websocket`: WebSocket gateway in Go
  - `chat-message-service`: Message storage and delivery in Go
- **AWS Services:** DynamoDB, SQS

---

## Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Node.js](https://nodejs.org/) (for frontend development)
- [AWS Account](https://aws.amazon.com/)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) (for cloud deployment)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) & [Helm](https://helm.sh/) (for Kubernetes deployment)
- [Terraform](https://www.terraform.io/) (for infrastructure provisioning)

---

## Local Development

### 1. Configure AWS Credentials

The backend services require access to AWS DynamoDB and SQS. Set your AWS credentials in an `.env` file at the project root:

```
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
AWS_REGION=your-region
```

You can copy `.env.example` and fill in your values.

### 2. Start Backend Services

Run the backend services using Docker Compose:

```fish
# From the project root
docker compose up --build
```

This will start:

- `gateway-websocket` (port 8080)
- `chat-message-service` (port 8081)

### 3. Start the Frontend

In a new terminal:

```fish
cd chatapp-frontend
npm install
npm run dev
```

The frontend will be available at [http://localhost:5173](http://localhost:5173) by default.

---

## Environment Variables

- Frontend: see `chatapp-frontend/.env.example` for required variables (e.g., `VITE_WS_URL`, `VITE_CHAT_SERVICE_URL`).
- Backend: requires AWS credentials as above.

---

## Cloud Deployment

### 1. Provision Infrastructure

Use Terraform scripts in `deploy/terraform/` to provision AWS resources (EKS, DynamoDB, SQS, etc.):

```fish
cd deploy/terraform
terraform init
terraform apply
```

### 2. Build and Push Docker Images

Authenticate with AWS ECR and push images:

```fish
make push
```

### 3. Deploy to Kubernetes

Use Helm charts in `deploy/helm/`:

```fish
make helm-deploy
```

---

## Step-by-Step Cloud Deployment Guide

### 1. Prerequisites

- AWS account and CLI configured (`aws configure`)
- Docker & Docker Compose
- Node.js (for frontend build)
- kubectl, Helm, and Terraform installed

### 2. Provision AWS Infrastructure

```fish
cd deploy/terraform
terraform init
terraform apply
```

- This creates EKS, DynamoDB, SQS, VPC, S3 and other resources.
- Note the outputs: SQS queue URL, DynamoDB table name, EKS cluster name and S3 bucket name.

### 3. Configure kubectl for EKS

```fish
aws eks --region <region> update-kubeconfig --name <eks_cluster_name>
```

- Use the EKS cluster name from Terraform output.

### 4. Build and Push Docker Images

```fish
make push
```

- This builds and pushes both backend images to ECR.

### 5. Create Kubernetes Secrets

```fish
make create-secrets
```

- This command creates secrets in Kubernetes from your `.env` file (with AWS keys, SQS URL, DynamoDB table, etc).

### 6. Deploy Backend Services with Helm

```fish
make helm-deploy
```

- This deploys both backend services to EKS using Helm charts.

### 7. Deploy ServiceAccount (if not automated)

```fish
kubectl apply -f k8s/serviceaccount-chatapp.yaml
```

### 8. Deploy/Configure Frontend

- Build the frontend:
  ```fish
  cd chatapp-frontend
  npm install
  npm run build
  ```
- Deploy the static files to S3 (bucket is provisioned by Terraform):
  ```fish
  aws s3 sync chatapp-frontend/dist s3://<your-s3-bucket-name> --delete
  # Or use the output from Terraform: frontend_static_site_bucket
  ```
- (Optional) Set up CloudFront for CDN delivery.
- Set `VITE_WS_URL` and `VITE_CHAT_SERVICE_URL` to the correct endpoints.

#### S3 Bucket Details

- The S3 bucket for the frontend static site is created by Terraform (`deploy/terraform/s3.tf`).
- The bucket name is controlled by the `s3_bucket_name` variable (default: `chatapp-frontend-bucket`).
- The bucket name is output as `frontend_static_site_bucket` after `terraform apply`.

**Tip:** In your CI/CD pipeline, set the `S3_BUCKET_NAME` secret to match the Terraform output for seamless deployment.

### 9. Test and Monitor

- Check `/health` endpoints for all services.
- Test chat functionality end-to-end.
- Monitor logs and set up alerts as needed.

---

## Project Structure

```
chatapp/
  chatapp-frontend/      # React frontend
  services/
    gateway-websocket/   # WebSocket backend (Go)
    chat-message-service/# Message backend (Go)
  deploy/
    helm/                # Helm charts
    terraform/           # Terraform scripts
  docker-compose.yaml    # Local dev orchestration
  Makefile               # Build & deploy helpers
```

---

## Notes

- Always start backend services before the frontend.
- AWS credentials are required for both local and cloud deployments.
- For production, review and secure your environment variables and secrets.

---

## License

MIT
