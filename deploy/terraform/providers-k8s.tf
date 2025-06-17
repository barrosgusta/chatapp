provider "kubernetes" {
  host                   = aws_eks_cluster.chatapp.endpoint
  cluster_ca_certificate = base64decode(aws_eks_cluster.chatapp.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.chatapp.token
}
