# VPC：自分専用のネットワーク（10.0.x.x）
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true # VPC 内で名前解決できる（RDS の接続先名を引くのに必要）
  enable_dns_hostnames = true # EC2 に DNS 名を付ける（VPC エンドポイント等で必要）

  tags = {
    Name = "fitness-app-vpc"
  }
}

# パブリックサブネット（1a）：NAT を置く
resource "aws_subnet" "public_1a" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = "ap-northeast-1a"

  tags = {
    Name = "fitness-app-public-1a"
  }
}

# パブリックサブネット（1c）：今は空。将来用
resource "aws_subnet" "public_1c" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "ap-northeast-1c"

  tags = {
    Name = "fitness-app-public-1c"
  }
}

# プライベートサブネット（1a）：Lambda・RDS
resource "aws_subnet" "private_1a" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.10.0/24"
  availability_zone = "ap-northeast-1a"

  tags = {
    Name = "fitness-app-private-1a"
  }
}

# プライベートサブネット（1c）：Lambda・RDS（RDS は AZ が 2 つ必要）
resource "aws_subnet" "private_1c" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.11.0/24"
  availability_zone = "ap-northeast-1c"

  tags = {
    Name = "fitness-app-private-1c"
  }
}

# インターネットゲートウェイ：VPC とインターネットの出入り口
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "fitness-app-igw"
  }
}

# パブリック用ルートテーブル（看板）
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "fitness-app-public-rt"
  }
}

# パブリックの看板に「どこ宛てでも IGW へ」を書く
resource "aws_route" "public_internet" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"                  # 宛先：全部
  gateway_id             = aws_internet_gateway.main.id # 送り先：IGW
}

# public_1a にパブリックの看板を立てる
resource "aws_route_table_association" "public_1a" {
  subnet_id      = aws_subnet.public_1a.id
  route_table_id = aws_route_table.public.id
}

# public_1c にパブリックの看板を立てる
resource "aws_route_table_association" "public_1c" {
  subnet_id      = aws_subnet.public_1c.id
  route_table_id = aws_route_table.public.id
}

# プライベート用ルートテーブル（看板）
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "fitness-app-private-rt"
  }
}

# private_1a にプライベートの看板を立てる
resource "aws_route_table_association" "private_1a" {
  subnet_id      = aws_subnet.private_1a.id
  route_table_id = aws_route_table.private.id
}

# private_1c にプライベートの看板を立てる
resource "aws_route_table_association" "private_1c" {
  subnet_id      = aws_subnet.private_1c.id
  route_table_id = aws_route_table.private.id
}

# default SG のルールを全部消す（指定し忘れても通さない）
resource "aws_default_security_group" "default" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "fitness-app-default-sg-unused"
  }
}
