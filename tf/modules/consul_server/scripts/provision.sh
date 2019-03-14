#!/bin/bash
set -ex

# Set hostname gets confusing otherwise on large servers
sudo hostname "$(cat ~/hostname)"
grep "$(cat ~/hostname)" /etc/hosts || echo "127.0.0.1     $(cat ~/hostname)" | sudo tee -a /etc/hosts
chmod 644 ~/hostname
sudo chown root:root ~/hostname
sudo mv ~/hostname /etc/hostname

# Install Consul
sudo chmod +x consul
sudo chown root:root consul
sudo mv consul /usr/local/bin

# Create consul user + data dir
getent passwd consul || sudo useradd --system --home /etc/consul.d --shell /bin/false consul
sudo mkdir --parents /opt/consul
sudo chown --recursive consul:consul /opt/consul

# Move systemd unit into place
sudo chown root:root consul.service
sudo mv consul.service /etc/systemd/system/consul.service

# Create configs
sudo mkdir --parents /etc/consul.d
sudo chown root:root consul.hcl
sudo mv consul.hcl /etc/consul.d/consul.hcl
sudo chown --recursive consul:consul /etc/consul.d
sudo chmod 640 /etc/consul.d/consul.hcl

# Launch cluster
sudo systemctl enable consul
sudo systemctl start consul
