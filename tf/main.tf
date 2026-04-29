provider "aws" {
  region = var.region
}

resource "aws_security_group" "mc" {
  name = "mc-sg"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.operator_cidr]
  }

  ingress {
    from_port   = 25565
    to_port     = 25565
    protocol    = "tcp"
    cidr_blocks = var.allowed_cidrs
  }

  ingress {
    from_port   = 24454
    to_port     = 24454
    protocol    = "udp"
    cidr_blocks = var.allowed_cidrs
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "mc" {
  ami                  = data.aws_ami.ubuntu.id
  instance_type        = var.instance_type
  key_name             = var.key_name
  iam_instance_profile = aws_iam_instance_profile.mc.name

  vpc_security_group_ids = [aws_security_group.mc.id]

  metadata_options {
    http_tokens   = "required"
    http_endpoint = "enabled"
  }

  user_data = templatefile("setup.sh.tftpl", {
    forge_url    = var.forge_url
    forge_sha256 = var.forge_sha256
  })

  tags = {
    Name    = "minecraft-server"
    Profile = var.profile_name
  }
}

output "instance_ip" {
  value = aws_instance.mc.public_ip
}

output "instance_id" {
  value = aws_instance.mc.id
}
