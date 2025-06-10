locals {
  services = [
    "gateway-websocket",
    "chat-message-service"
  ]
}

resource "aws_ecr_repository" "services" {
  for_each = toset(local.services)
  name     = "${var.project_prefix}-${each.key}"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "${var.project_prefix}-ecr-${each.key}"
  }
}
