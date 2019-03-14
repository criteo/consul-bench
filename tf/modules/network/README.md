## Network

This module is respondisble for creating a simple network for a consul datacenter(s) to live within. Since it is meant to be
isolated, this creates a network with public IPs, securing the instances via only whitelisting any access to the specified
CIDR ranges provided by the user.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| whitelisted\_cidrs | List of CIDRs to allow to access cluster | list | n/a | yes |
| cidr | Specifies cidr block the VPC will have for its network | string | `"172.21.0.0/16"` | no |

## Outputs

| Name | Description |
|------|-------------|
| access\_sg | Allows access for ssh to the specific whitelisted cidrs |
| subnet\_id | Aws ID of subnet created |
| vpc\_sg | Security group provided by module that allows intra-vpc communication on all nodes in the security group |

