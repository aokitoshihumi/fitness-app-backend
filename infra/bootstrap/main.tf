# state用のS3バケットだけを作る
terraform {
  required_version = "~> 1.17"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"

  default_tags {
    tags = {
      Project   = "fitness-app"
      ManagedBy = "terraform"
    }
  }
}

# ログイン中のアカウント ID を取る（バケット名に使う）
data "aws_caller_identity" "current" {}

# state を置くバケット
resource "aws_s3_bucket" "tfstate" {
  bucket = "fitness-app-tfstate-${data.aws_caller_identity.current.account_id}"

  # うっかり destroy しても消えない
  lifecycle {
    prevent_destroy = true
  }
}

# 上書き前の state も残す（壊しても戻せる）
resource "aws_s3_bucket_versioning" "tfstate" {
  bucket = aws_s3_bucket.tfstate.id

  versioning_configuration {
    status = "Enabled"
  }
}

# 保存時に暗号化する
resource "aws_s3_bucket_server_side_encryption_configuration" "tfstate" {
  bucket = aws_s3_bucket.tfstate.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# 公開アクセスを全部ブロック
resource "aws_s3_bucket_public_access_block" "tfstate" {
  bucket = aws_s3_bucket.tfstate.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# 作ったバケット名を表示する（prod の backend に書き写す）
output "bucket_name" {
  value = aws_s3_bucket.tfstate.bucket
}