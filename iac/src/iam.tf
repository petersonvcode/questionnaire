resource "aws_iam_role" "backend" {
  name = "q-backend-instance-role-${var.environment}"
  path = "/"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_policy" "backend" {
  name        = "q-backend-permissions-${var.environment}"
  description = "Permissions for the questionnaire backend instance in ${var.environment}"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = "ssm:GetParameter"
        Effect   = "Allow"
        Resource = [
          aws_ssm_parameter.backend_configuration.arn,
          aws_ssm_parameter.github_token.arn
        ]
      },
      {
        Action   = ["ec2:Describe*"],
        Effect   = "Allow",
        Resource = "*"
      },
      {
        Action = [
          "ec2:AuthorizeSecurityGroupIngress",
          "ec2:RevokeSecurityGroupIngress",
        ]
        Effect   = "Allow"
        Resource = aws_security_group.backend.arn
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "backend" {
  role       = aws_iam_role.backend.name
  policy_arn = aws_iam_policy.backend.arn
}