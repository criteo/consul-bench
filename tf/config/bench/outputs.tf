output "server_ips" {
  value = ["${module.consul_servers.public_ips}"]
}

output "bench_ip" {
  value = "${element(module.consul_bench.public_ips, 0)}"
}
