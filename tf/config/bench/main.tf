terraform {
  require_version = "~> 0.11"
}

provider "aws" {
  version = "~> 1.57"
  region  = "${var.region}"
}

provider "local" {
  version = "~> 1.1"
}

provider "template" {
  version = "~> 2.0"
}

provider "null" {
  version = "~> 2.0"
}

locals {
  allowed_cidrs = [
    "91.199.242.236/32",
  ]
}

module "network" {
  source = "../../modules/network"

  whitelisted_cidrs = ["${local.allowed_cidrs}"]
}

module "consul_servers" {
  source = "../../modules/consul_server"

  subnet_id     = "${module.network.subnet_id}"
  node_count    = 1
  consul_path   = "${var.consul_path}"
  ssh_key       = "${var.ssh_key}"
  ssh_key_id    = "${var.ssh_key_id}"
  custom_config = "${var.custom_config}"

  sg_ids = [
    "${module.network.vpc_sg}",
    "${module.network.access_sg}",
  ]
}


module "consul_bench" {
  source = "../../modules/consul_bench"

  subnet_id         = "${module.network.subnet_id}"
  ssh_key           = "${var.ssh_key}"
  ssh_key_id        = "${var.ssh_key_id}"
  consul_bench_path = "${var.consul_bench_path}"
  target            = "${module.consul_servers.private_ips}"

  sg_ids = [
    "${module.network.vpc_sg}",
    "${module.network.access_sg}",
  ]
}
