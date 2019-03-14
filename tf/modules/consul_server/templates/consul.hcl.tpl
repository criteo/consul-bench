# Some dummy info
datacenter = "${dc}"

performance {
  raft_multiplier = 1
}

# Enable it as a server
server = ${server}
${expected}

# Configure other servers to join
retry_join = [
${nodes}
]

