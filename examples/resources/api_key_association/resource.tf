resource "searchstax_api_key_association" "example" {
  account_name   = "my_account"
  api_key        = var.api_key
  deployment_uid = "ss123456"
}

variable "api_key" {
  type      = string
  sensitive = true
}
