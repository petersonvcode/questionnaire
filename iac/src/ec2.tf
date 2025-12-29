locals {
  userdata_files = [
    "${path.module}/scripts/00-get-config.sh", # Needs to be the first in list
    "${path.module}/scripts/10-install-backend-dependencies.sh",
    "${path.module}/scripts/20-setup-backend.sh",
    "${path.module}/scripts/99-main.sh", # Needs to be the last in list
  ]
  userdata_content = join("\n", [for path in local.userdata_files : file(path)])
  userdata_gzip    = base64gzip(local.userdata_content)

  aws_ebs_volume_zone = "${var.region}a"
}

resource "aws_instance" "backend" {
  ami                  = data.aws_ami.amazon2.id
  instance_type        = "t3.nano"
  availability_zone    = local.aws_ebs_volume_zone
  iam_instance_profile = aws_iam_instance_profile.backend.name

  user_data_base64            = local.userdata_gzip
  user_data_replace_on_change = true

  associate_public_ip_address = true
  key_name                    = aws_key_pair.backend.key_name

  security_groups = [
    aws_security_group.allow_ssh.name,
    aws_security_group.backend.name
  ]

  tags = {
    Name        = "Questionnaire Backend - ${title(var.environment)}"
    environment = var.environment
  }

  metadata_options {
    instance_metadata_tags = "enabled"
  }

  depends_on = [ aws_s3_bucket.server_assets ]
}

data "aws_ami" "amazon2" {
  most_recent = true

  filter {
    name   = "name"
    values = ["al2023-ami-2023.9.20251117.1-kernel-6.1-x86_64"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["137112412989"]
}

resource "aws_iam_instance_profile" "backend" {
  name = "q-backend-instance-profile-${var.environment}"
  role = aws_iam_role.backend.name
}

resource "aws_eip" "backend" {
  instance = aws_instance.backend.id

  tags = {
    Name = "q-backend-eip-${var.environment}"
  }
}

resource "aws_key_pair" "backend" {
  key_name   = "q-backend"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAII/I0Zkg0vW7DF7vmECSceoc0sp0OgdJCfYzby5acbae petersonvgama@gmail.com"
}

resource "aws_security_group" "allow_ssh" {
  name        = "q-allow-ssh-${var.environment}"
  description = "Allows ssh from dev machine"

  vpc_id = data.aws_vpc.default.id
}

resource "aws_security_group_rule" "allow_ssh_ingress" {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = ["187.121.9.51/32"]
  security_group_id = aws_security_group.allow_ssh.id
  description       = "Allow SSH from devs machines"

  lifecycle {
    ignore_changes = [cidr_blocks]
  }
}

resource "aws_security_group_rule" "allow_vps_access_to_backend" {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = ["172.62.106.70/32"]
  security_group_id = aws_security_group.allow_ssh.id
  description       = "Allow SSH from VPS machine. Required for automated deployments."
}

resource "aws_security_group" "backend" {
  name        = "q-backend-${var.environment}"
  description = "Allows backend to communicate with the world"

  vpc_id = data.aws_vpc.default.id

  tags = {
    # Id tag is used by the backend startup script
    Id = "q-backend-main-sg"
  }
}

resource "aws_security_group_rule" "allow_backend_port" {
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.backend.id
}

resource "aws_security_group_rule" "allow_all_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.backend.id
  description       = "Allow all outbound traffic"
}