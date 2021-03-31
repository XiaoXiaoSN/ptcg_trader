locals {
  cluster_name                 = "ptcg.10oz.tw"
  master_autoscaling_group_ids = [aws_autoscaling_group.master-ap-northeast-1a-masters-ptcg-10oz-tw.id]
  master_security_group_ids    = [aws_security_group.masters-ptcg-10oz-tw.id]
  masters_role_arn             = aws_iam_role.masters-ptcg-10oz-tw.arn
  masters_role_name            = aws_iam_role.masters-ptcg-10oz-tw.name
  node_autoscaling_group_ids   = [aws_autoscaling_group.nodes-ap-northeast-1a-ptcg-10oz-tw.id]
  node_security_group_ids      = [aws_security_group.nodes-ptcg-10oz-tw.id]
  node_subnet_ids              = [aws_subnet.ap-northeast-1a-ptcg-10oz-tw.id]
  nodes_role_arn               = aws_iam_role.nodes-ptcg-10oz-tw.arn
  nodes_role_name              = aws_iam_role.nodes-ptcg-10oz-tw.name
  region                       = "ap-northeast-1"
  route_table_public_id        = aws_route_table.ptcg-10oz-tw.id
  subnet_ap-northeast-1a_id    = aws_subnet.ap-northeast-1a-ptcg-10oz-tw.id
  vpc_cidr_block               = aws_vpc.ptcg-10oz-tw.cidr_block
  vpc_id                       = aws_vpc.ptcg-10oz-tw.id
}

output "cluster_name" {
  value = "ptcg.10oz.tw"
}

output "master_autoscaling_group_ids" {
  value = [aws_autoscaling_group.master-ap-northeast-1a-masters-ptcg-10oz-tw.id]
}

output "master_security_group_ids" {
  value = [aws_security_group.masters-ptcg-10oz-tw.id]
}

output "masters_role_arn" {
  value = aws_iam_role.masters-ptcg-10oz-tw.arn
}

output "masters_role_name" {
  value = aws_iam_role.masters-ptcg-10oz-tw.name
}

output "node_autoscaling_group_ids" {
  value = [aws_autoscaling_group.nodes-ap-northeast-1a-ptcg-10oz-tw.id]
}

output "node_security_group_ids" {
  value = [aws_security_group.nodes-ptcg-10oz-tw.id]
}

output "node_subnet_ids" {
  value = [aws_subnet.ap-northeast-1a-ptcg-10oz-tw.id]
}

output "nodes_role_arn" {
  value = aws_iam_role.nodes-ptcg-10oz-tw.arn
}

output "nodes_role_name" {
  value = aws_iam_role.nodes-ptcg-10oz-tw.name
}

output "region" {
  value = "ap-northeast-1"
}

output "route_table_public_id" {
  value = aws_route_table.ptcg-10oz-tw.id
}

output "subnet_ap-northeast-1a_id" {
  value = aws_subnet.ap-northeast-1a-ptcg-10oz-tw.id
}

output "vpc_cidr_block" {
  value = aws_vpc.ptcg-10oz-tw.cidr_block
}

output "vpc_id" {
  value = aws_vpc.ptcg-10oz-tw.id
}

# provider "aws" {
#   region = "ap-northeast-1"
# }

resource "aws_autoscaling_group" "master-ap-northeast-1a-masters-ptcg-10oz-tw" {
  enabled_metrics = ["GroupDesiredCapacity", "GroupInServiceInstances", "GroupMaxSize", "GroupMinSize", "GroupPendingInstances", "GroupStandbyInstances", "GroupTerminatingInstances", "GroupTotalInstances"]
  launch_template {
    id      = aws_launch_template.master-ap-northeast-1a-masters-ptcg-10oz-tw.id
    version = aws_launch_template.master-ap-northeast-1a-masters-ptcg-10oz-tw.latest_version
  }
  max_size            = 1
  metrics_granularity = "1Minute"
  min_size            = 1
  name                = "master-ap-northeast-1a.masters.ptcg.10oz.tw"
  tag {
    key                 = "KubernetesCluster"
    propagate_at_launch = true
    value               = "ptcg.10oz.tw"
  }
  tag {
    key                 = "Name"
    propagate_at_launch = true
    value               = "master-ap-northeast-1a.masters.ptcg.10oz.tw"
  }
  tag {
    key                 = "k8s.io/cluster-autoscaler/node-template/label/kops.k8s.io/instancegroup"
    propagate_at_launch = true
    value               = "master-ap-northeast-1a"
  }
  tag {
    key                 = "k8s.io/cluster-autoscaler/node-template/label/kubernetes.io/role"
    propagate_at_launch = true
    value               = "master"
  }
  tag {
    key                 = "k8s.io/cluster-autoscaler/node-template/label/node-role.kubernetes.io/master"
    propagate_at_launch = true
    value               = ""
  }
  tag {
    key                 = "k8s.io/role/master"
    propagate_at_launch = true
    value               = "1"
  }
  tag {
    key                 = "kops.k8s.io/instancegroup"
    propagate_at_launch = true
    value               = "master-ap-northeast-1a"
  }
  tag {
    key                 = "kubernetes.io/cluster/ptcg.10oz.tw"
    propagate_at_launch = true
    value               = "owned"
  }
  vpc_zone_identifier = [aws_subnet.ap-northeast-1a-ptcg-10oz-tw.id]
}

resource "aws_autoscaling_group" "nodes-ap-northeast-1a-ptcg-10oz-tw" {
  enabled_metrics = ["GroupDesiredCapacity", "GroupInServiceInstances", "GroupMaxSize", "GroupMinSize", "GroupPendingInstances", "GroupStandbyInstances", "GroupTerminatingInstances", "GroupTotalInstances"]
  launch_template {
    id      = aws_launch_template.nodes-ap-northeast-1a-ptcg-10oz-tw.id
    version = aws_launch_template.nodes-ap-northeast-1a-ptcg-10oz-tw.latest_version
  }
  max_size            = 1
  metrics_granularity = "1Minute"
  min_size            = 1
  name                = "nodes-ap-northeast-1a.ptcg.10oz.tw"
  tag {
    key                 = "KubernetesCluster"
    propagate_at_launch = true
    value               = "ptcg.10oz.tw"
  }
  tag {
    key                 = "Name"
    propagate_at_launch = true
    value               = "nodes-ap-northeast-1a.ptcg.10oz.tw"
  }
  tag {
    key                 = "k8s.io/cluster-autoscaler/node-template/label/kops.k8s.io/instancegroup"
    propagate_at_launch = true
    value               = "nodes-ap-northeast-1a"
  }
  tag {
    key                 = "k8s.io/cluster-autoscaler/node-template/label/kubernetes.io/role"
    propagate_at_launch = true
    value               = "node"
  }
  tag {
    key                 = "k8s.io/cluster-autoscaler/node-template/label/node-role.kubernetes.io/node"
    propagate_at_launch = true
    value               = ""
  }
  tag {
    key                 = "k8s.io/role/node"
    propagate_at_launch = true
    value               = "1"
  }
  tag {
    key                 = "kops.k8s.io/instancegroup"
    propagate_at_launch = true
    value               = "nodes-ap-northeast-1a"
  }
  tag {
    key                 = "kubernetes.io/cluster/ptcg.10oz.tw"
    propagate_at_launch = true
    value               = "owned"
  }
  vpc_zone_identifier = [aws_subnet.ap-northeast-1a-ptcg-10oz-tw.id]
}

resource "aws_ebs_volume" "a-etcd-events-ptcg-10oz-tw" {
  availability_zone = "ap-northeast-1a"
  encrypted         = false
  size              = 20
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "a.etcd-events.ptcg.10oz.tw"
    "k8s.io/etcd/events"                 = "a/a"
    "k8s.io/role/master"                 = "1"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
  }
  type = "gp2"
}

resource "aws_ebs_volume" "a-etcd-main-ptcg-10oz-tw" {
  availability_zone = "ap-northeast-1a"
  encrypted         = false
  size              = 20
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "a.etcd-main.ptcg.10oz.tw"
    "k8s.io/etcd/main"                   = "a/a"
    "k8s.io/role/master"                 = "1"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
  }
  type = "gp2"
}

resource "aws_iam_instance_profile" "masters-ptcg-10oz-tw" {
  name = "masters.ptcg.10oz.tw"
  role = aws_iam_role.masters-ptcg-10oz-tw.name
}

resource "aws_iam_instance_profile" "nodes-ptcg-10oz-tw" {
  name = "nodes.ptcg.10oz.tw"
  role = aws_iam_role.nodes-ptcg-10oz-tw.name
}

resource "aws_iam_role_policy" "masters-ptcg-10oz-tw" {
  name   = "masters.ptcg.10oz.tw"
  policy = file("${path.module}/data/aws_iam_role_policy_masters.ptcg.10oz.tw_policy")
  role   = aws_iam_role.masters-ptcg-10oz-tw.name
}

resource "aws_iam_role_policy" "nodes-ptcg-10oz-tw" {
  name   = "nodes.ptcg.10oz.tw"
  policy = file("${path.module}/data/aws_iam_role_policy_nodes.ptcg.10oz.tw_policy")
  role   = aws_iam_role.nodes-ptcg-10oz-tw.name
}

resource "aws_iam_role" "masters-ptcg-10oz-tw" {
  assume_role_policy = file("${path.module}/data/aws_iam_role_masters.ptcg.10oz.tw_policy")
  name               = "masters.ptcg.10oz.tw"
}

resource "aws_iam_role" "nodes-ptcg-10oz-tw" {
  assume_role_policy = file("${path.module}/data/aws_iam_role_nodes.ptcg.10oz.tw_policy")
  name               = "nodes.ptcg.10oz.tw"
}

resource "aws_internet_gateway" "ptcg-10oz-tw" {
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "ptcg.10oz.tw"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
  }
  vpc_id = aws_vpc.ptcg-10oz-tw.id
}

resource "aws_key_pair" "kubernetes-ptcg-10oz-tw-7c18a4652564f54aebc6f68873b4194c" {
  key_name   = "kubernetes.ptcg.10oz.tw-7c:18:a4:65:25:64:f5:4a:eb:c6:f6:88:73:b4:19:4c"
  public_key = file("${path.module}/data/aws_key_pair_kubernetes.ptcg.10oz.tw-7c18a4652564f54aebc6f68873b4194c_public_key")
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "ptcg.10oz.tw"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
  }
}

resource "aws_launch_template" "master-ap-northeast-1a-masters-ptcg-10oz-tw" {
  block_device_mappings {
    device_name = "/dev/sda1"
    ebs {
      delete_on_termination = true
      encrypted             = false
      volume_size           = 64
      volume_type           = "gp2"
    }
  }
  iam_instance_profile {
    name = aws_iam_instance_profile.masters-ptcg-10oz-tw.id
  }
  image_id      = "ami-0f08dd35fe25614c9"
  instance_type = "t3.medium"
  key_name      = aws_key_pair.kubernetes-ptcg-10oz-tw-7c18a4652564f54aebc6f68873b4194c.id
  lifecycle {
    create_before_destroy = true
  }
  metadata_options {
    http_endpoint               = "enabled"
    http_put_response_hop_limit = 1
    http_tokens                 = "optional"
  }
  name = "master-ap-northeast-1a.masters.ptcg.10oz.tw"
  network_interfaces {
    associate_public_ip_address = true
    delete_on_termination       = true
    security_groups             = [aws_security_group.masters-ptcg-10oz-tw.id]
  }
  tag_specifications {
    resource_type = "instance"
    tags = {
      "KubernetesCluster"                                                            = "ptcg.10oz.tw"
      "Name"                                                                         = "master-ap-northeast-1a.masters.ptcg.10oz.tw"
      "k8s.io/cluster-autoscaler/node-template/label/kops.k8s.io/instancegroup"      = "master-ap-northeast-1a"
      "k8s.io/cluster-autoscaler/node-template/label/kubernetes.io/role"             = "master"
      "k8s.io/cluster-autoscaler/node-template/label/node-role.kubernetes.io/master" = ""
      "k8s.io/role/master"                                                           = "1"
      "kops.k8s.io/instancegroup"                                                    = "master-ap-northeast-1a"
      "kubernetes.io/cluster/ptcg.10oz.tw"                                           = "owned"
    }
  }
  tag_specifications {
    resource_type = "volume"
    tags = {
      "KubernetesCluster"                                                            = "ptcg.10oz.tw"
      "Name"                                                                         = "master-ap-northeast-1a.masters.ptcg.10oz.tw"
      "k8s.io/cluster-autoscaler/node-template/label/kops.k8s.io/instancegroup"      = "master-ap-northeast-1a"
      "k8s.io/cluster-autoscaler/node-template/label/kubernetes.io/role"             = "master"
      "k8s.io/cluster-autoscaler/node-template/label/node-role.kubernetes.io/master" = ""
      "k8s.io/role/master"                                                           = "1"
      "kops.k8s.io/instancegroup"                                                    = "master-ap-northeast-1a"
      "kubernetes.io/cluster/ptcg.10oz.tw"                                           = "owned"
    }
  }
  tags = {
    "KubernetesCluster"                                                            = "ptcg.10oz.tw"
    "Name"                                                                         = "master-ap-northeast-1a.masters.ptcg.10oz.tw"
    "k8s.io/cluster-autoscaler/node-template/label/kops.k8s.io/instancegroup"      = "master-ap-northeast-1a"
    "k8s.io/cluster-autoscaler/node-template/label/kubernetes.io/role"             = "master"
    "k8s.io/cluster-autoscaler/node-template/label/node-role.kubernetes.io/master" = ""
    "k8s.io/role/master"                                                           = "1"
    "kops.k8s.io/instancegroup"                                                    = "master-ap-northeast-1a"
    "kubernetes.io/cluster/ptcg.10oz.tw"                                           = "owned"
  }
  user_data = filebase64("${path.module}/data/aws_launch_template_master-ap-northeast-1a.masters.ptcg.10oz.tw_user_data")
}

resource "aws_launch_template" "nodes-ap-northeast-1a-ptcg-10oz-tw" {
  block_device_mappings {
    device_name = "/dev/sda1"
    ebs {
      delete_on_termination = true
      encrypted             = false
      volume_size           = 128
      volume_type           = "gp2"
    }
  }
  iam_instance_profile {
    name = aws_iam_instance_profile.nodes-ptcg-10oz-tw.id
  }
  image_id      = "ami-0f08dd35fe25614c9"
  instance_type = "t3.medium"
  key_name      = aws_key_pair.kubernetes-ptcg-10oz-tw-7c18a4652564f54aebc6f68873b4194c.id
  lifecycle {
    create_before_destroy = true
  }
  metadata_options {
    http_endpoint               = "enabled"
    http_put_response_hop_limit = 1
    http_tokens                 = "optional"
  }
  name = "nodes-ap-northeast-1a.ptcg.10oz.tw"
  network_interfaces {
    associate_public_ip_address = true
    delete_on_termination       = true
    security_groups             = [aws_security_group.nodes-ptcg-10oz-tw.id]
  }
  tag_specifications {
    resource_type = "instance"
    tags = {
      "KubernetesCluster"                                                          = "ptcg.10oz.tw"
      "Name"                                                                       = "nodes-ap-northeast-1a.ptcg.10oz.tw"
      "k8s.io/cluster-autoscaler/node-template/label/kops.k8s.io/instancegroup"    = "nodes-ap-northeast-1a"
      "k8s.io/cluster-autoscaler/node-template/label/kubernetes.io/role"           = "node"
      "k8s.io/cluster-autoscaler/node-template/label/node-role.kubernetes.io/node" = ""
      "k8s.io/role/node"                                                           = "1"
      "kops.k8s.io/instancegroup"                                                  = "nodes-ap-northeast-1a"
      "kubernetes.io/cluster/ptcg.10oz.tw"                                         = "owned"
    }
  }
  tag_specifications {
    resource_type = "volume"
    tags = {
      "KubernetesCluster"                                                          = "ptcg.10oz.tw"
      "Name"                                                                       = "nodes-ap-northeast-1a.ptcg.10oz.tw"
      "k8s.io/cluster-autoscaler/node-template/label/kops.k8s.io/instancegroup"    = "nodes-ap-northeast-1a"
      "k8s.io/cluster-autoscaler/node-template/label/kubernetes.io/role"           = "node"
      "k8s.io/cluster-autoscaler/node-template/label/node-role.kubernetes.io/node" = ""
      "k8s.io/role/node"                                                           = "1"
      "kops.k8s.io/instancegroup"                                                  = "nodes-ap-northeast-1a"
      "kubernetes.io/cluster/ptcg.10oz.tw"                                         = "owned"
    }
  }
  tags = {
    "KubernetesCluster"                                                          = "ptcg.10oz.tw"
    "Name"                                                                       = "nodes-ap-northeast-1a.ptcg.10oz.tw"
    "k8s.io/cluster-autoscaler/node-template/label/kops.k8s.io/instancegroup"    = "nodes-ap-northeast-1a"
    "k8s.io/cluster-autoscaler/node-template/label/kubernetes.io/role"           = "node"
    "k8s.io/cluster-autoscaler/node-template/label/node-role.kubernetes.io/node" = ""
    "k8s.io/role/node"                                                           = "1"
    "kops.k8s.io/instancegroup"                                                  = "nodes-ap-northeast-1a"
    "kubernetes.io/cluster/ptcg.10oz.tw"                                         = "owned"
  }
  user_data = filebase64("${path.module}/data/aws_launch_template_nodes-ap-northeast-1a.ptcg.10oz.tw_user_data")
}

resource "aws_route_table_association" "ap-northeast-1a-ptcg-10oz-tw" {
  route_table_id = aws_route_table.ptcg-10oz-tw.id
  subnet_id      = aws_subnet.ap-northeast-1a-ptcg-10oz-tw.id
}

resource "aws_route_table" "ptcg-10oz-tw" {
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "ptcg.10oz.tw"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
    "kubernetes.io/kops/role"            = "public"
  }
  vpc_id = aws_vpc.ptcg-10oz-tw.id
}

resource "aws_route" "route-0-0-0-0--0" {
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.ptcg-10oz-tw.id
  route_table_id         = aws_route_table.ptcg-10oz-tw.id
}

resource "aws_security_group_rule" "all-master-to-master" {
  from_port                = 0
  protocol                 = "-1"
  security_group_id        = aws_security_group.masters-ptcg-10oz-tw.id
  source_security_group_id = aws_security_group.masters-ptcg-10oz-tw.id
  to_port                  = 0
  type                     = "ingress"
}

resource "aws_security_group_rule" "all-master-to-node" {
  from_port                = 0
  protocol                 = "-1"
  security_group_id        = aws_security_group.nodes-ptcg-10oz-tw.id
  source_security_group_id = aws_security_group.masters-ptcg-10oz-tw.id
  to_port                  = 0
  type                     = "ingress"
}

resource "aws_security_group_rule" "all-node-to-node" {
  from_port                = 0
  protocol                 = "-1"
  security_group_id        = aws_security_group.nodes-ptcg-10oz-tw.id
  source_security_group_id = aws_security_group.nodes-ptcg-10oz-tw.id
  to_port                  = 0
  type                     = "ingress"
}

resource "aws_security_group_rule" "https-external-to-master-0-0-0-0--0" {
  cidr_blocks       = ["0.0.0.0/0"]
  from_port         = 443
  protocol          = "tcp"
  security_group_id = aws_security_group.masters-ptcg-10oz-tw.id
  to_port           = 443
  type              = "ingress"
}

resource "aws_security_group_rule" "master-egress" {
  cidr_blocks       = ["0.0.0.0/0"]
  from_port         = 0
  protocol          = "-1"
  security_group_id = aws_security_group.masters-ptcg-10oz-tw.id
  to_port           = 0
  type              = "egress"
}

resource "aws_security_group_rule" "node-egress" {
  cidr_blocks       = ["0.0.0.0/0"]
  from_port         = 0
  protocol          = "-1"
  security_group_id = aws_security_group.nodes-ptcg-10oz-tw.id
  to_port           = 0
  type              = "egress"
}

resource "aws_security_group_rule" "node-to-master-tcp-1-2379" {
  from_port                = 1
  protocol                 = "tcp"
  security_group_id        = aws_security_group.masters-ptcg-10oz-tw.id
  source_security_group_id = aws_security_group.nodes-ptcg-10oz-tw.id
  to_port                  = 2379
  type                     = "ingress"
}

resource "aws_security_group_rule" "node-to-master-tcp-2382-4000" {
  from_port                = 2382
  protocol                 = "tcp"
  security_group_id        = aws_security_group.masters-ptcg-10oz-tw.id
  source_security_group_id = aws_security_group.nodes-ptcg-10oz-tw.id
  to_port                  = 4000
  type                     = "ingress"
}

resource "aws_security_group_rule" "node-to-master-tcp-4003-65535" {
  from_port                = 4003
  protocol                 = "tcp"
  security_group_id        = aws_security_group.masters-ptcg-10oz-tw.id
  source_security_group_id = aws_security_group.nodes-ptcg-10oz-tw.id
  to_port                  = 65535
  type                     = "ingress"
}

resource "aws_security_group_rule" "node-to-master-udp-1-65535" {
  from_port                = 1
  protocol                 = "udp"
  security_group_id        = aws_security_group.masters-ptcg-10oz-tw.id
  source_security_group_id = aws_security_group.nodes-ptcg-10oz-tw.id
  to_port                  = 65535
  type                     = "ingress"
}

resource "aws_security_group_rule" "ssh-external-to-master-0-0-0-0--0" {
  cidr_blocks       = ["0.0.0.0/0"]
  from_port         = 22
  protocol          = "tcp"
  security_group_id = aws_security_group.masters-ptcg-10oz-tw.id
  to_port           = 22
  type              = "ingress"
}

resource "aws_security_group_rule" "ssh-external-to-node-0-0-0-0--0" {
  cidr_blocks       = ["0.0.0.0/0"]
  from_port         = 22
  protocol          = "tcp"
  security_group_id = aws_security_group.nodes-ptcg-10oz-tw.id
  to_port           = 22
  type              = "ingress"
}

resource "aws_security_group" "masters-ptcg-10oz-tw" {
  description = "Security group for masters"
  name        = "masters.ptcg.10oz.tw"
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "masters.ptcg.10oz.tw"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
  }
  vpc_id = aws_vpc.ptcg-10oz-tw.id
}

resource "aws_security_group" "nodes-ptcg-10oz-tw" {
  description = "Security group for nodes"
  name        = "nodes.ptcg.10oz.tw"
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "nodes.ptcg.10oz.tw"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
  }
  vpc_id = aws_vpc.ptcg-10oz-tw.id
}

resource "aws_subnet" "ap-northeast-1a-ptcg-10oz-tw" {
  availability_zone = "ap-northeast-1a"
  cidr_block        = "172.20.32.0/19"
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "ap-northeast-1a.ptcg.10oz.tw"
    "SubnetType"                         = "Public"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
    "kubernetes.io/role/elb"             = "1"
  }
  vpc_id = aws_vpc.ptcg-10oz-tw.id
}

resource "aws_vpc_dhcp_options_association" "ptcg-10oz-tw" {
  dhcp_options_id = aws_vpc_dhcp_options.ptcg-10oz-tw.id
  vpc_id          = aws_vpc.ptcg-10oz-tw.id
}

resource "aws_vpc_dhcp_options" "ptcg-10oz-tw" {
  domain_name         = "ap-northeast-1.compute.internal"
  domain_name_servers = ["AmazonProvidedDNS"]
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "ptcg.10oz.tw"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
  }
}

resource "aws_vpc" "ptcg-10oz-tw" {
  cidr_block           = "172.20.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    "KubernetesCluster"                  = "ptcg.10oz.tw"
    "Name"                               = "ptcg.10oz.tw"
    "kubernetes.io/cluster/ptcg.10oz.tw" = "owned"
  }
}

terraform {
  required_version = ">= 0.12.0"
}
