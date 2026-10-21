#!/usr/bin/env bash
# ==============================================================================
# Intili Sandbox — Packer AMI Provisioning Script
# ==============================================================================
# This script provisions an Ubuntu 24.04 LTS instance to bake a production-ready,
# zero-download AMI for Intili Sandbox.
#
# It installs:
#   - Docker Engine & Docker Compose Plugin (Official Docker Repositories)
#   - Build tools (git, make, curl, jq, net-tools)
#   - Copies the Packer-uploaded opi-p4-sandbox source to /opt/opi-p4-sandbox
#   - Pre-builds and caches all Docker images locally inside the AMI
#   - Configures and enables /etc/systemd/system/opi-sandbox.service
# ==============================================================================

set -euo pipefail

export DEBIAN_FRONTEND=noninteractive

echo "=========================================================="
echo "🚀 [1/6] Updating APT packages and installing dependencies..."
echo "=========================================================="
apt-get update -y
apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    git \
    make \
    jq \
    net-tools \
    tar \
    software-properties-common

echo "=========================================================="
echo "🐳 [2/6] Installing official Docker CE & Docker Compose..."
echo "=========================================================="
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update -y
apt-get install -y --no-install-recommends \
    docker-ce \
    docker-ce-cli \
    containerd.io \
    docker-buildx-plugin \
    docker-compose-plugin

# Ensure docker daemon runs and start on boot
systemctl enable docker
systemctl start docker

# Add standard cloud users to docker group
usermod -aG docker ubuntu || true
if id "ec2-user" &>/dev/null; then
    usermod -aG docker ec2-user
fi

echo "=========================================================="
echo "📂 [3/6] Setting up /opt/opi-p4-sandbox repository..."
echo "=========================================================="
mkdir -p /opt/opi-p4-sandbox

if [ ! -f "/tmp/opi-p4-sandbox/docker-compose.yml" ]; then
    echo "Packer source upload is missing or incomplete."
    exit 1
fi
cp -a /tmp/opi-p4-sandbox/. /opt/opi-p4-sandbox/
rm -rf /tmp/opi-p4-sandbox

# Set ownership
chown -R ubuntu:ubuntu /opt/opi-p4-sandbox || true

echo "=========================================================="
echo "🔨 [4/6] Pre-building Docker images into AMI storage..."
echo "=========================================================="
cd /opt/opi-p4-sandbox

# Build images directly so EC2 instances boot instantly with zero download
docker compose build --pull

# Verify images are stored in local Docker cache
docker images

echo "=========================================================="
echo "⚙️  [5/6] Configuring systemd auto-start service..."
echo "=========================================================="
if [ -f "/tmp/opi-sandbox.service" ]; then
    cp /tmp/opi-sandbox.service /etc/systemd/system/opi-sandbox.service
else
    cp /opt/opi-p4-sandbox/packer/opi-sandbox.service /etc/systemd/system/opi-sandbox.service
fi

chmod 644 /etc/systemd/system/opi-sandbox.service
systemctl daemon-reload
systemctl enable opi-sandbox.service

echo "=========================================================="
echo "🧹 [6/6] Cleaning up temporary files for AMI release..."
echo "=========================================================="
# Clean apt caches
apt-get clean
rm -rf /var/lib/apt/lists/*

# Clean cloud-init and history
rm -rf /tmp/* /var/tmp/*
truncate -s 0 /etc/machine-id || true
rm -f /var/lib/dbus/machine-id || true
history -c || true

echo "=========================================================="
echo "✅ AMI Baking completed successfully!"
echo "   Intili Sandbox is ready for instant boot."
echo "=========================================================="
