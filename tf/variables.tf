variable "key_name" {
  default = "minecraft"
}

variable "region" {
  default = "sa-east-1"
}

variable "instance_type" {
  default = "m6i.large"
}

variable "operator_cidr" {
  type = string
}

variable "allowed_cidrs" {
  type = list(string)
}

variable "forge_url" {
  type = string
}

variable "forge_sha256" {
  type = string
}

variable "profile_name" {
  type = string
}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }
}
