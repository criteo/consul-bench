## Consul Server

Used to create consul servers. Can either come up in server mode, or if an IP or list of IPs are provided will come up as consul
agents for an existing consul cluster

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| consul\_path | Path on local filesystem to a consul binary to use for the system | string | n/a | yes |
| node\_count | Number of consul servers in the cluster (should be an odd number for consensus) | string | n/a | yes |
| sg\_ids | List of security groups to apply to the server | list | n/a | yes |
| ssh\_key | Private key used to authenticate to the node | string | n/a | yes |
| ssh\_key\_id | Name of key stored in AWS for autenticating with the node | string | n/a | yes |
| subnet\_id | ID for the subnet to create the servers within | string | n/a | yes |
| custom\_config | Extra configuration to add to the cluster, [required: json] | string | `"{}"` | no |
| dc\_name | Specifies the name of the datacenter, important if used in multi-dc use cases | string | `"dc0"` | no |
| instance\_type | Type of instance to use when creating servers | string | `"t2.large"` | no |
| root\_block\_size | Specifies the size (in GB) of the root block device | string | `"30"` | no |
| secondary\_block\_name | Specifies the name which the secondary block device will be mounted as | string | `"xvdb"` | no |
| secondary\_block\_size | Specifies the size (in GB) of the second block device | string | `"30"` | no |
| server\_nodes | List of node IPs to connect to as an agent, leave empty and nodes created will become a server cluster | list | `<list>` | no |

## Outputs

| Name | Description |
|------|-------------|
| private\_ips | Private IP addresses of the instances created |
| public\_ips | Public IP addresses of instances created |

