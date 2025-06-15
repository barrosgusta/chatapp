resource "aws_s3_bucket" "frontend_static_site" {
  bucket = var.s3_bucket_name

  tags = {
    Name        = "${var.project_prefix}-frontend-static-site"
    Environment = "dev"
  }
}
