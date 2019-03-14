locals {
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-*-18.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "bench" {
  ami                    = "${data.aws_ami.ubuntu.id}"
  instance_type          = "${var.instance_type}"
  subnet_id              = "${var.subnet_id}"
  vpc_security_group_ids = ["${var.sg_ids}"]
  key_name               = "${var.ssh_key_id}"

  tags {
    Name = "TF - consul-bench ${count.index}"
    Host = "${format("consul-bench%03d", count.index)}"
  }

  count = 1

  root_block_device {
    volume_size = "${var.root_block_size}"
  }

  ebs_block_device {
    device_name = "${var.secondary_block_name}"
    volume_size = "${var.secondary_block_size}"
  }
}

resource "null_resource" "provision_consul_bench" {
  triggers {
    uuid = "${uuid()}"
  }

  connection {
    host        = "${element(aws_instance.bench.*.public_ip, 0)}"
    type        = "ssh"
    user        = "ubuntu"
    private_key = "${file(var.ssh_key)}"
  }

  provisioner "file" {
    source      = "${var.consul_bench_path}"
    destination = "~/consul-bench"
  }

  provisioner "file" {
    source = "${path.module}/scripts/provision.sh"
    destination = "~/provision.sh"
  }

  provisioner "remote-exec" {
    inline = [
        "chmod +x ~/provision.sh",
        "~/provision.sh ${element(var.target, 0)} > stdout.log 2> stderr.log"
    ]
  }
}