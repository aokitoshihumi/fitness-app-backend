# DB をおくサブネットのくみ（2つのAZが必要らしい）
resource "aws_db_subnet_group" "main" {
  name       = "fitness-app-db-subnet-group"
  subnet_ids = [aws_subnet.private_1a.id, aws_subnet.private_1c.id]

  tags = {
    Name = "fitness-app-db-subnet-group"
  }
}

# DB のパスワード（記号なし・32文字）
resource "random_password" "db" {
  length  = 32
  special = false # 接続文字列が壊れるかもなので記号false
}

# RDS 本体
resource "aws_db_instance" "main" {
  identifier     = "fitness-app-db"
  engine         = "postgres"
  engine_version = "18" # メジャー指定
  instance_class = "db.t4g.micro"

  allocated_storage = 20    # ディスクの容量
  storage_type      = "gp3" # SSDだがgp3は速くて安い
  storage_encrypted = true  # 暗号化

  db_name  = "fitness"                 # 最初に作られるDBの名前
  username = "fitness"                 # DBの管理ユーザ名
  password = random_password.db.result # 上で作った32文字のパスワード

  db_subnet_group_name   = aws_db_subnet_group.main.name # 置き場所（private1a, private1c）
  vpc_security_group_ids = [aws_security_group.rds.id]   # rds-sg(Lambdaの5432だけ通す)
  publicly_accessible    = false                         # インターネットから見えなくする
  multi_az               = false                         # スタンバイは作らない

  backup_retention_period   = 7                      # 自動バックアップを7日分残す
  deletion_protection       = true                   # 削除できなくする
  skip_final_snapshot       = false                  # 消す時最後にスナップショットを撮る
  final_snapshot_identifier = "fitness-app-db-final" # そのスナップショットの名前

}