resource "searchstax_tags_set" "example" {
  account_name   = "my_account"
  deployment_uid = "ss123456"
  tags           = ["production", "terraform"]
}
