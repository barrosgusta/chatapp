variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1"
}

variable "project_prefix" {
  description = "Prefix for all resource names"
  type        = string
  default     = "chatapp"
}

variable "dynamodb_chat_messages_table" {
  description = "DynamoDB table for chat messages"
  type        = string
  default     = "chat_messages"
}

variable "sqs_queue_name" {
  description = "SQS queue for processing messages"
  type        = string
  default     = "message-task-queue"
}

variable "vpc_cidr" {
  description = "CIDR for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "s3_bucket_name" {
  description = "S3 bucket name for frontend static site deployment"
  type        = string
  default     = "chatapp-frontend-bucket"
}