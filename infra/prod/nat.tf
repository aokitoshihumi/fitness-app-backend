##
## NATは高いのでt4g.nanoを使いたい
##

# NATの固定パブリックIP（外からみた送信元）
resource "aws_eip" "nat" {
  domain = "vpc"

  tags = {
    Name = "fitness-app-nat-eip"
  }
}

# fck-nat: プライベートサブネットの外向き通信を中継する
module "fck_nat" {
  source  = "RaJiska/fck-nat/aws"
  version = "1.6.1"

  name          = "fitness-app-nat"
  vpc_id        = aws_vpc.main.id
  subnet_id     = aws_subnet.public_1a.id
  instance_type = "t4g.nano"

  # HA= 高可用性をtrueにすることでASGに属させる
  ha_mode            = true                        # 落ちたら自動で作り直す
  eip_allocation_ids = [aws_eip.nat.allocation_id] # 上の固定 IP を使う
  attach_ssm_policy  = true                        # 踏み台で使う（管理用の入口）をきめる

  # プライベートの看板に「外向きは NAT へ」を足す
  update_route_tables = true
  route_tables_ids = {
    private = aws_route_table.private.id
  }
}
