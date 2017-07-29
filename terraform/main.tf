variable project {
  default = "json-charlieegan3"
}

terraform {
  backend "s3" {
    bucket = "charlieegan3-www-terraform-state"
    region = "us-east-1"
    key    = "json-charlieegan3.tfstate"
  }
}

resource "aws_iam_user" "default" {
  name          = "json-charlieegan3"
  force_destroy = true
}

resource "aws_iam_user_policy" "default" {
  name = "s3"
  user = "${aws_iam_user.default.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:*"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::charlieegan3-www-website-content/status.json"
    },
    {
      "Action": [
        "s3:ListBucket"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::charlieegan3-www-website-content"
    },
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": [
        "cloudfront:CreateInvalidation"
      ]
    }
  ]
}
EOF
}
