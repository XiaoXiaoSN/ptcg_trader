variable "region" {
  default = "ap-northeast-1"
}

# terraform {
#   backend "remote" {
#     organization = "10oz"

#     workspaces {
#       name = "PTCG-Trader"
#     }
#   }
# }

locals {
  environment             = "Dev"
  kops_state_bucket_name  = "ptcg-kops-state"
  kubernetes_cluster_name = "ptcg.10oz.tw"
  ingress_ips             = ["10.0.0.100/32", "10.0.0.101/32"]
  vpc_name                = "ptcg-vpc"

  tags = {
    environment = "Dev"
    terraform   = true
  }
}
