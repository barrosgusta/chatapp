resource "aws_dynamodb_table" "chat_messages" {
  name           = var.dynamodb_chat_messages_table
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "Id"

  attribute {
    name = "Id"
    type = "S"
  }

  tags = {
    Name        = "${var.project_prefix}-dynamodb-chat"
    Environment = "dev"
  }
}