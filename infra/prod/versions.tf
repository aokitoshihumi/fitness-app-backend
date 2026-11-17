# 動かすための前提条件（go.mod の役割）
terraform {
  # terraformのバージョン指定（1.17 以上 2.0 未満）
  required_version = "~> 1.17"

  # AWS と話すための部品（プロバイダー）
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0" # 6.x に固定
    }
  }
}
