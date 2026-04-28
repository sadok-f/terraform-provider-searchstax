resource "searchstax_dns_record" "example" {
  account_name = var.account_name
  name         = "myalias"
  deployment   = var.deployment_uid
  ttl          = "300"
}

