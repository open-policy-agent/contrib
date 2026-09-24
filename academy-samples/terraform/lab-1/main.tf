resource "aws_instance" "app_server" {
  ami           = "ami-830c94e3"
  instance_type = "t2.micro"
  root_block_device {
    volume_size = 8
  }
}

resource "aws_ebs_volume" "volume1" {
  availability_zone = "us-west-2a"
  size              = 2
}