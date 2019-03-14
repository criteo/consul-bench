variable "ssh_key" {
  description = "Path to ssh key, used for provisioning servers"
  type        = "string"
}

variable "ssh_key_id" {
  description = "Name of key in aws"
  type        = "string"
}

variable "region" {
  description = "Region in AWS to build cluster within"
  type        = "string"
  default     = "us-east-1"
}

variable "dc_name" {
  description = "sets the name of the datacenter"
  type        = "string"
  default     = "dc0"
}

variable "consul_path" {
  description = "Path to consul binary to install on agent nodes"
  type        = "string"
  default     = "./consul"
}

variable "consul_bench_path" {
  description = "Path to consul binary to install on agent nodes"
  type        = "string"
  default     = "./consul-bench"
}


variable "custom_config" {
  description = "Path to a JSON file with configuration for the consul cluster"
  type        = "string"
  default     = "./files/empty.json"
}
