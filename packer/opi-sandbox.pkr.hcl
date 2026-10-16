# ==============================================================================
# OPI & P4 Developer Sandbox — HashiCorp Packer Template (AWS AMI)
# ==============================================================================
# Builds a pre-baked, production-ready Amazon Machine Image (AMI) for AWS EC2.
# When an instance is launched from this AMI:
#   1. Zero downloads required at boot time.
#   2. Docker & Docker Compose are pre-installed.
#   3. OPI API Server & P4 BMv2 images are pre-baked in Docker image cache.
#   4. systemd service auto-starts the sandbox on boot.
#
# Usage:
#   packer init packer/opi-sandbox.pkr.hcl
#   packer build -var="aws_region=us-east-1" packer/opi-sandbox.pkr.hcl
# ==============================================================================

packer {
  required_version = ">= 1.9.0"
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

# --- Variables ---

variable "aws_region" {
  type        = string
  default     = "us-east-1"
  description = "AWS Region to build and publish the AMI"
}

variable "instance_type" {
  type        = string
  default     = "t3.medium"
  description = "EC2 instance type used during AMI baking"
}

variable "ami_prefix" {
  type        = string
  default     = "opi-p4-sandbox"
  description = "Prefix for the resulting AMI name"
}

variable "ami_version" {
  type        = string
  default     = "v1.0.0"
  description = "Version tag for the sandbox AMI"
}

variable "vpc_id" {
  type        = string
  default     = ""
  description = "Optional VPC ID to build the AMI inside"
}

variable "subnet_id" {
  type        = string
  default     = ""
  description = "Optional Subnet ID to build the AMI inside"
}

# --- Locals ---

locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

# --- Builder Source ---

source "amazon-ebs" "ubuntu" {
  region        = var.aws_region
  instance_type = var.instance_type
  ami_name      = "${var.ami_prefix}-${var.ami_version}-${local.timestamp}"
  ami_description = "Pre-baked OPI and P4 developer sandbox with Docker, BMv2 switch, and Web Dashboard"

  # Canonical Ubuntu 24.04 LTS (Noble Numbat) 64-bit (x86_64)
  source_ami_filter {
    filters = {
      "name"                = "ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-*"
      "root-device-type"    = "ebs"
      "virtualization-type" = "hvm"
    }
    owners      = ["099720109477"] # Canonical Official AWS Account ID
    most_recent = true
  }

  ssh_username = "ubuntu"
  vpc_id       = var.vpc_id != "" ? var.vpc_id : null
  subnet_id    = var.subnet_id != "" ? var.subnet_id : null

  # 30 GB gp3 root volume for Docker images and build caches
  launch_block_device_mappings {
    device_name           = "/dev/sda1"
    volume_size           = 30
    volume_type           = "gp3"
    iops                  = 3000
    throughput            = 125
    delete_on_termination = true
  }

  tags = {
    Name        = "${var.ami_prefix}-${var.ami_version}"
    Version     = var.ami_version
    Project     = "OPI-P4-Sandbox"
    Environment = "Development"
    OS          = "Ubuntu 24.04 LTS"
    BuiltBy     = "Packer"
  }
}

# --- Build Definition ---

build {
  name    = "opi-sandbox-ami"
  sources = ["source.amazon-ebs.ubuntu"]

  # Upload systemd service file
  provisioner "file" {
    source      = "${path.root}/opi-sandbox.service"
    destination = "/tmp/opi-sandbox.service"
  }

  # The file provisioner requires a pre-existing destination directory.
  provisioner "shell" {
    inline = ["mkdir -p /tmp/opi-p4-sandbox"]
  }

  # Bake the exact source tree checked out by the caller or CI workflow.
  provisioner "file" {
    source      = "${path.root}/../"
    destination = "/tmp/opi-p4-sandbox"
  }

  # Execute provisioning script as sudo
  provisioner "shell" {
    script          = "${path.root}/setup.sh"
    execute_command = "sudo -E bash '{{ .Path }}'"
  }

  # Final verification post-provision
  provisioner "shell" {
    inline = [
      "echo '=== Docker Version ==='",
      "docker --version",
      "docker compose version",
      "echo '=== Pre-Baked Images ==='",
      "docker images",
      "echo '=== Systemd Auto-Boot Service ==='",
      "systemctl is-enabled opi-sandbox.service || true",
      "echo '=== Verification Complete ==='"
    ]
  }
}
