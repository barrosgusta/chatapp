data "aws_eks_cluster" "chatapp" {
  name = aws_eks_cluster.chatapp.name
}

data "aws_eks_cluster_auth" "chatapp" {
  name = aws_eks_cluster.chatapp.name
}

data "tls_certificate" "eks" {
  url = data.aws_eks_cluster.chatapp.identity[0].oidc[0].issuer
}

resource "aws_iam_openid_connect_provider" "eks" {
  url             = data.aws_eks_cluster.chatapp.identity[0].oidc[0].issuer
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = [data.tls_certificate.eks.certificates[0].sha1_fingerprint]
}

resource "aws_iam_role" "service_role" {
  name = "${var.project_prefix}-service-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = {
        Federated = aws_iam_openid_connect_provider.eks.arn
      },
      Action = "sts:AssumeRoleWithWebIdentity",
      Condition = {
        StringEquals = {
          "${replace(data.aws_eks_cluster.chatapp.identity[0].oidc[0].issuer, "https://", "")}:sub" = "system:serviceaccount:default:chatapp-serviceaccount"
        }
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "service_policy_attach" {
  role       = aws_iam_role.service_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_role_policy_attachment" "service_policy_sqs" {
  role       = aws_iam_role.service_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_iam_role_policy_attachment" "service_policy_dynamodb" {
  role       = aws_iam_role.service_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
}

