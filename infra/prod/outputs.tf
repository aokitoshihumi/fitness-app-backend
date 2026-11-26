output "vpc_id" {
  description = "VPC の ID"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "パブリックサブネットの ID [1a, 1c] の順"
  value       = [aws_subnet.public_1a.id, aws_subnet.public_1c.id]
}

output "private_subnet_ids" {
  description = "プライベートサブネットの ID [1a, 1c] の順"
  value       = [aws_subnet.private_1a.id, aws_subnet.private_1c.id]
}

output "nat_public_ip" {
  description = "NAT の固定パブリック IP"
  value       = aws_eip.nat.public_ip
}

output "nat_autoscaling_group_arn" {
  description = "NAT インスタンスを管理する Auto Scaling グループの ARN"
  value       = module.fck_nat.autoscaling_group_arn
}

output "db_address" {
  description = "RDSの接続先（ホスト名）"
  value       = aws_db_instance.main.address


}
output "db_password" {
  description = "DB のパスワード（terraform output -raw db_password で見る）"
  value       = random_password.db.result
  sensitive   = true
}