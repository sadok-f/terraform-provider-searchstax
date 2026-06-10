resource "searchstax_deployment_user" "example" {
  account_name   = "my_account"
  deployment_uid = "ss123456"
  username       = "solruser"
  password       = var.solr_user_password
  role           = "Admin"
}

variable "solr_user_password" {
  type      = string
  sensitive = true
}
