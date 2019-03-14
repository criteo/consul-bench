variable "subnet_id" {
  description = "ID for the subnet to create the servers within"
  type        = "string"
}

variable "sg_ids" {
  description = "List of security groups to apply to the server"
  type        = "list"
}

variable "node_count" {
  description = "Number of consul servers in the cluster (should be an odd number for consensus)"
  type        = "string"
}

variable "consul_path" {
  description = "Path on local filesystem to a consul binary to use for the system"
  type        = "string"
}

variable "ssh_key" {
  description = "Private key used to authenticate to the node"
  type        = "string"
}

variable "ssh_key_id" {
  description = "Name of key stored in AWS for autenticating with the node"
  type        = "string"
}

variable "server_nodes" {
  description = "List of node IPs to connect to as an agent, leave empty and nodes created will become a server cluster"
  type        = "list"
  default     = []
}

variable "instance_type" {
  description = "Type of instance to use when creating servers"
  type        = "string"
  default     = "t2.large"
}

variable "root_block_size" {
  description = "Specifies the size (in GB) of the root block device"
  type        = "string"
  default     = "30"
}

variable "secondary_block_size" {
  description = "Specifies the size (in GB) of the second block device"
  type        = "string"
  default     = "30"
}

variable "secondary_block_name" {
  description = "Specifies the name which the secondary block device will be mounted as"
  type        = "string"
  default     = "xvdb"
}

variable "dc_name" {
  description = "Specifies the name of the datacenter, important if used in multi-dc use cases"
  type        = "string"
  default     = "dc0"
}

variable "custom_config" {
  description = "Extra configuration to add to the cluster, [required: json]"
  type        = "string"
  default     = "{}"
}
