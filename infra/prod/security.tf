# lambdaのセキュリティグループ
resource "aws_security_group" "lambda" {
  name        = "fitness-app-lambda-sg"
  description = "Security group for Lambda"
  vpc_id      = aws_vpc.main.id

  tags = {
    Name = "fitness-app-lambda-sg"
  }
}

# RDS 用のセキュリティグループ
resource "aws_security_group" "rds" {
  name        = "fitness-app-rds-sg"
  description = "Security group for RDS"
  vpc_id      = aws_vpc.main.id

  tags = {
    Name = "fitness-app-rds-sg"
  }
}

# LambdaからRDSへのルール
resource "aws_vpc_security_group_egress_rule" "lambda_to_rds" {
  security_group_id            = aws_security_group.lambda.id # lambda-sgに付けるルール
  description                  = "Lambda to RDS"
  ip_protocol                  = "tcp"
  from_port                    = 5432
  to_port                      = 5432
  referenced_security_group_id = aws_security_group.rds.id # lambdaはRDS宛なら通信を通す
}

# Lambdaから外部のHTTPS（Firebase宛の通信とか）
resource "aws_vpc_security_group_egress_rule" "lambda_to_https" {
  security_group_id = aws_security_group.lambda.id
  description       = "HTTPS to internet"
  ip_protocol       = "tcp"
  from_port         = 443
  to_port           = 443
  cidr_ipv4         = "0.0.0.0/0" # 宛先のIPは問わない
}

# RDS <- Lambdaだけを受け取る
resource "aws_vpc_security_group_ingress_rule" "rds_from_lambda" {
  security_group_id            = aws_security_group.rds.id
  description                  = "PostgreSQL from Lambda"
  ip_protocol                  = "tcp"
  from_port                    = 5432
  to_port                      = 5432
  referenced_security_group_id = aws_security_group.lambda.id # IP指定じゃなくてもSG指定で通信できちゃう
}