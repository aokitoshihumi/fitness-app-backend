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

  # state の置き場所（bucket は backend.hcl から渡す）
  backend "s3" {
    key          = "prod/terraform.tfstate" # バケット内のファイル名
    region       = "ap-northeast-1"
    encrypt      = true
    use_lockfile = true # 同時に apply しないためのロック
  }
}
