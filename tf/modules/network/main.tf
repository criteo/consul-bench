resource "aws_vpc" "net" {
  cidr_block = "${var.cidr}"
}

resource "aws_subnet" "s" {
  cidr_block              = "${aws_vpc.net.cidr_block}"
  vpc_id                  = "${aws_vpc.net.id}"
  map_public_ip_on_launch = true
}

resource "aws_internet_gateway" "default" {
  vpc_id = "${aws_vpc.net.id}"
}

resource "aws_route" "internet" {
  route_table_id         = "${aws_vpc.net.main_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.default.id}"
}

resource "aws_security_group" "access" {
  name        = "access-sg"
  description = "allows remote access to specified address"
  vpc_id      = "${aws_vpc.net.id}"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["${var.whitelisted_cidrs}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "vpc_all" {
  name        = "vpc-allow-all-sg"
  description = "allows all nodes within vpc in this security group to communicate"
  vpc_id      = "${aws_vpc.net.id}"

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = -1
    self      = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}
