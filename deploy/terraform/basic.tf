provider "aws" {
  #   access_key = "$AWS_ACCESS_KEY_ID"
  #   secret_key = "$AWS_SECRET_ACCESS_KEY"

  # kops Note: S3 requires --create-bucket-configuration LocationConstraint=<region> 
  # for regions other than us-east-1.
  region = var.region
}

resource "aws_s3_bucket" "ptcg-bucket" {
  bucket = "ptcg-bucket-tf"
  acl    = "private"

  tags = {
    Name        = "ptcg bucket"
    Environment = "Dev"
  }
}

resource "aws_route53_zone" "dev" {
  name = "ptcg.10oz.tw"
}
