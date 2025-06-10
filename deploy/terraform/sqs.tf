resource "aws_sqs_queue" "message_queue" {
  name = var.sqs_queue_name

  tags = {
    Name        = "${var.project_prefix}-sqs-message-queue"
    Environment = "dev"
  }
}
