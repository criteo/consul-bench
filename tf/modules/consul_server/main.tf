locals {
  server                = "${length(var.server_nodes) == 0 ? true : false}"
  formated_server_nodes = "${join(",\n", formatlist("  \"%v\"", var.server_nodes))}"
  formated_aws_nodes    = "${join(",\n", formatlist("  \"%v\"", aws_instance.server.*.private_ip))}"
  retry_join            = "${local.server == true ? local.formated_aws_nodes : local.formated_server_nodes}"
  expected              = "${local.server == true ? format("bootstrap_expect = %v", length(aws_instance.server.*.private_ip)) : ""}"
  name                  = "${local.server == true ? "server" : "agent"}"
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

data "template_file" "hostname" {
  template = "${file("${path.module}/templates/empty.tpl")}"

  vars {
    content = "${format("consul-%v%03d-%v", local.name, count.index, var.dc_name)}"
  }

  count = "${var.node_count}"
}

data "template_file" "consul_config" {
  template = "${file("${path.module}/templates/consul.hcl.tpl")}"

  vars {
    nodes    = "${local.retry_join}"
    expected = "${local.expected}"
    server   = "${local.server}"
    dc       = "${var.dc_name}"
  }
}

resource "aws_instance" "server" {
  ami                    = "${data.aws_ami.ubuntu.id}"
  instance_type          = "${var.instance_type}"
  subnet_id              = "${var.subnet_id}"
  vpc_security_group_ids = ["${var.sg_ids}"]
  key_name               = "${var.ssh_key_id}"

  tags {
    Name = "TF - consul-${local.name} ${count.index}"
    Host = "${format("consul-${local.name}%03d", count.index)}"
  }

  count = "${var.node_count}"

  root_block_device {
    volume_size = "${var.root_block_size}"
  }

  ebs_block_device {
    device_name = "${var.secondary_block_name}"
    volume_size = "${var.secondary_block_size}"
  }
}

resource "null_resource" "provision_consul" {
  count = "${var.node_count}"

  triggers {
    uuid = "${uuid()}"
  }

  connection {
    host        = "${element(aws_instance.server.*.public_ip, count.index)}"
    type        = "ssh"
    user        = "ubuntu"
    private_key = "${file(var.ssh_key)}"
  }

  provisioner "file" {
    source      = "${var.consul_path}"
    destination = "~/consul"
  }

  provisioner "file" {
    content     = "${element(data.template_file.hostname.*.rendered, count.index)}"
    destination = "~/hostname"
  }

  provisioner "file" {
    source      = "${path.module}/files/consul.service"
    destination = "~/consul.service"
  }

  provisioner "file" {
    content     = "${data.template_file.consul_config.rendered}"
    destination = "~/consul.hcl"
  }

  provisioner "file" {
    content     = "${var.custom_config}"
    destination = "~/extra_config.json"
  }

  provisioner "file" {
    source = "${path.module}/scripts/provision.sh"
    destination = "~/provision.sh"
  }

  provisioner "remote-exec" {
    inline = [
        "chmod +x ~/provision.sh",
        "~/provision.sh > stdout.log 2> stderr.log"
    ]
  }
}
