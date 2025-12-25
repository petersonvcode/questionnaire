resource "aws_ebs_volume" "backend_persistence" {
  availability_zone = local.aws_ebs_volume_zone
  size              = 20
  type              = "gp3"

  tags = {
    Name = "q-backend-persistence-${var.environment}"
  }
}

resource "aws_volume_attachment" "backend_persistence" {
  device_name = "/dev/xvdf"
  volume_id   = aws_ebs_volume.backend_persistence.id
  instance_id = aws_instance.backend.id
}