variable "key_name" {
  default = "minecraft"
}

variable "region" {
  default = "sa-east-1"
}

variable "instance_type" {
  default = "m7i.xlarge"
}

variable "operator_cidr" {
  type = string
}

variable "allowed_cidrs" {
  type = list(string)
}

variable "loader" {
  type    = string
  default = "forge"
}

variable "minecraft_version" {
  type    = string
  default = ""
}

variable "fabric_loader_version" {
  type    = string
  default = ""
}

variable "forge_url" {
  type    = string
  default = ""
}

variable "forge_sha256" {
  type    = string
  default = ""
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
