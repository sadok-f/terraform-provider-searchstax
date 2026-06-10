data "searchstax_api_key_deployments" "example" {
  account_name = "my_account"
  api_key      = var.api_key
}

variable "api_key" {
  type      = string
  sensitive = true
}
