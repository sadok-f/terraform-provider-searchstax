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

data "searchstax_private_vpc" "example" {
  account_name = "example_account"
}


output "example_s_private_vpc" {
  value = data.searchstax_private_vpc.example
}