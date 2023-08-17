terraform {
  required_providers {
    searchstax = {
      source = "hashicorp.com/sadok-f/searchstax"
    }
  }
}

provider "searchstax" {
  username = var.ssx_username
  password = var.ssx_pwd
}

data "searchstax_deployments" "example" {
  account_name = "example_account_name"
}


output "example_deployments" {
  value = data.searchstax_deployments.example
}