variable "cidr" {
  description = "Specifies cidr block the VPC will have for its network"
  type        = "string"
  default     = "172.21.0.0/16"
}

variable "whitelisted_cidrs" {
  description = "List of CIDRs to allow to access cluster"
  type        = "list"
}
