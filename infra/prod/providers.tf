# ログイン情報は書かない（export AWS_PROFILE=fitness で渡す）
provider "aws" {
  region = "ap-northeast-1" # 東京

  # 全リソースに自動で付けるタグ
  default_tags {
    tags = {
      Project   = "fitness-app"
      ManagedBy = "terraform" # 画面から手で変えない、の印
    }
  }
}
