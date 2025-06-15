output "sqs_queue_url" {
  description = "URL of the SQS message task queue"
  value       = aws_sqs_queue.message_queue.id
}

output "dynamodb_chat_messages_table" {
  description = "Name of the DynamoDB table for chat messages"
  value       = aws_dynamodb_table.chat_messages.name
}

output "eks_cluster_name" {
  description = "EKS Cluster name"
  value       = aws_eks_cluster.chatapp.name
}

output "frontend_static_site_bucket" {
  description = "Name of the S3 bucket for the frontend static site"
  value       = aws_s3_bucket.frontend_static_site.bucket
}