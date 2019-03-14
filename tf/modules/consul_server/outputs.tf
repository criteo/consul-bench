output "private_ips" {
  value       = ["${aws_instance.server.*.private_ip}"]
  description = "Private IP addresses of the instances created"
}

output "public_ips" {
  value       = ["${aws_instance.server.*.public_ip}"]
  description = "Public IP addresses of instances created"
}
