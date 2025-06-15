# EKS Worker Node Security Group
resource "aws_security_group" "eks_nodes" {
  name        = "${var.project_prefix}-eks-nodes-sg"
  description = "Security group for EKS worker nodes"
  vpc_id      = aws_vpc.main.id

  # Allow all traffic within the node group (for pod-to-pod communication)
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
  }

  # Allow inbound from EKS cluster security group (control plane)
  ingress {
    from_port                = 1025
    to_port                  = 65535
    protocol                 = "tcp"
    security_groups          = [aws_eks_cluster.chatapp.vpc_config[0].cluster_security_group_id]
    description              = "Allow EKS control plane to worker nodes"
  }
  ingress {
    from_port                = 443
    to_port                  = 443
    protocol                 = "tcp"
    security_groups          = [aws_eks_cluster.chatapp.vpc_config[0].cluster_security_group_id]
    description              = "Allow EKS control plane to worker nodes (HTTPS)"
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_prefix}-eks-nodes-sg"
  }
}
# Allow ALB to reach EKS nodes on app ports (using EKS cluster security group)
resource "aws_security_group_rule" "nodes_from_alb" {
  type                     = "ingress"
  from_port                = 8080
  to_port                  = 8081
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.alb.id
  security_group_id        = aws_eks_cluster.chatapp.vpc_config[0].cluster_security_group_id
}
# ALB Security Group
resource "aws_security_group" "alb" {
  name        = "${var.project_prefix}-alb-sg"
  description = "Allow HTTP and app ports to ALB"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8081
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Output ALB Security Group ID for use in annotations or reference
output "alb_security_group_id" {
  value = aws_security_group.alb.id
}
# VPC
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "${var.project_prefix}-vpc"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.project_prefix}-igw"
  }
}

# Public Subnets
resource "aws_subnet" "public" {
  count                   = 2
  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 4, count.index)
  map_public_ip_on_launch = true
  availability_zone       = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name = "${var.project_prefix}-public-${count.index}"
    "${format("kubernetes.io/cluster/%s-eks-cluster", var.project_prefix)}" = "shared"
    "kubernetes.io/role/elb" = "1"
  }
}

# Route Table
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "${var.project_prefix}-rt"
  }
}

# Route Table Association
resource "aws_route_table_association" "public" {
  count          = length(aws_subnet.public[*].id)
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Availability Zones
data "aws_availability_zones" "available" {
  state = "available"
}

# Private Subnets
resource "aws_subnet" "private" {
  count                   = 2
  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 4, count.index + 2)
  map_public_ip_on_launch = false
  availability_zone       = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name = "${var.project_prefix}-private-${count.index}"
    "${format("kubernetes.io/cluster/%s-eks-cluster", var.project_prefix)}" = "shared"
    "kubernetes.io/role/internal-elb" = "1"
  }
}

# Private Route Table
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.project_prefix}-private-rt"
  }
}

# Private Route Table Association
resource "aws_route_table_association" "private" {
  count          = length(aws_subnet.private[*].id)
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private.id
}

# Elastic IP for NAT Gateway
resource "aws_eip" "nat" {
  domain = "vpc"
}

# NAT Gateway in the first public subnet
resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public[0].id
  depends_on    = [aws_internet_gateway.igw]
}

# Route from private route table to NAT Gateway
resource "aws_route" "private_nat_gateway" {
  route_table_id         = aws_route_table.private.id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.nat.id
}