# Install AWS Load Balancer Controller via Helm
install-alb-controller:
	helm repo add eks https://aws.github.io/eks-charts
	helm repo update
	helm upgrade --install aws-load-balancer-controller eks/aws-load-balancer-controller \
	  -n kube-system \
	  --set clusterName=chatapp-eks-cluster \
	  --set serviceAccount.create=false \
	  --set serviceAccount.name=aws-load-balancer-controller \
	  --set region=us-east-1 \
	  --set vpcId=vpc-0f26a809bc43f1785
AWS_REGION=us-east-1
ACCOUNT_ID=767397757187
ECR_BASE=$(ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

SERVICES=gateway-websocket chat-message-service

.PHONY: login build tag push helm-deploy create-secrets

login:
	aws ecr get-login-password --region $(AWS_REGION) \
		| docker login --username AWS --password-stdin $(ECR_BASE)

build:
	@for service in $(SERVICES); do \
		echo "🔨 Building $$service..."; \
		docker build -t chatapp-$$service ./services/$$service; \
	done

tag:
	@for service in $(SERVICES); do \
		echo "🏷️  Tagging $$service..."; \
		docker tag chatapp-$$service:latest $(ECR_BASE)/chatapp-$$service:latest; \
	done

push: login build tag
	@for service in $(SERVICES); do \
		echo "🚀 Pushing $$service to ECR..."; \
		docker push $(ECR_BASE)/chatapp-$$service:latest; \
	done

helm-deploy:
	@for service in $(SERVICES); do \
		echo "📦 Deploying $$service via Helm..."; \
		helm upgrade --install chatapp-$$service ./deploy/helm/$$service \
		  --namespace chatapp --create-namespace \
		  --values ./deploy/helm/$$service/values.yaml; \
	done

create-secrets:
	@for service in $(SERVICES); do \
		echo "🔐 Criando secret para $$service..."; \
		kubectl create secret generic $$service-secrets \
			--from-env-file=.env \
			--namespace chatapp \
			--dry-run=client -o yaml | kubectl apply -f -; \
	done

apply-ingress:
	kubectl apply -f k8s/serviceaccount-chatapp.yaml
	kubectl apply -f k8s/ingress.yaml