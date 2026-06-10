terraform {
  required_providers {
    searchstax = {
      source = "hashicorp.com/sadok-f/searchstax"
    }
  }
}

provider "searchstax" {
  # Credentials can also be supplied via environment variables:
  # SEARCHSTAX_USERNAME, SEARCHSTAX_PASSWORD, SEARCHSTAX_HOST
  username = var.ssx_username
  password = var.ssx_pwd
  host     = var.ssx_host
}

variable "ssx_username" {
  type        = string
  description = "SearchStax account email."
}

variable "ssx_pwd" {
  type        = string
  sensitive   = true
  description = "SearchStax account password."
}

variable "ssx_host" {
  type        = string
  default     = "https://app.searchstax.com/api/rest/v2"
  description = "SearchStax REST API base URL."
}
