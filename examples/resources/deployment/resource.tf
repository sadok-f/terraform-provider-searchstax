terraform {
  required_providers {
    searchstax = {
      source = "hashicorp.com/sadok-f/searchstax"
    }
  }
}
variable "ssx_username" {
  type = string
}

variable "ssx_pwd" {
  type = string
}

variable "account_name" {
  type    = string
  default = "account_nameocp"
}

variable "private_vpc_name" {
  type    = string
  default = "account_name-na-ne1-vpc"
}

provider "searchstax" {
  username = var.ssx_username
  password = var.ssx_pwd
}

data "searchstax_private_vpc" "ssx_private_vpc" {
  account_name = var.account_name
}

locals {
  private_vpc_index = index(data.searchstax_private_vpc.ssx_private_vpc.private_vpc_list.*.name, var.private_vpc_name)
  private_vpc_id    = data.searchstax_private_vpc.ssx_private_vpc.private_vpc_list[local.private_vpc_index].id
}

resource "searchstax_deployment" "ssx_deployment" {
  account_name             = var.account_name
  name                     = "from-terraform-provider-tst-02"
  application              = "Solr"
  application_version      = "8.11.2"
  termination_lock         = "false"
  plan_type                = "DedicatedDeployment"
  plan                     = "NDC4-GCP-G"
  region_id                = "northamerica-northeast1"
  cloud_provider_id        = "gcp"
  num_additional_app_nodes = "0"
  private_vpc              = local.private_vpc_id
}


output "ssx_deployment_output" {
  value = searchstax_deployment.ssx_deployment.uid
}