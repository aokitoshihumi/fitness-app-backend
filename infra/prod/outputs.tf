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
