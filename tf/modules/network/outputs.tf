output "subnet_id" {
  value       = "${aws_subnet.s.id}"
  description = "Aws ID of subnet created"
}

output "vpc_sg" {
  value       = "${aws_security_group.vpc_all.id}"
  description = "Security group provided by module that allows intra-vpc communication on all nodes in the security group"
}

output "access_sg" {
  value       = "${aws_security_group.access.id}"
  description = "Allows access for ssh to the specific whitelisted cidrs"
}
