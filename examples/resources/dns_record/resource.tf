resource "searchstax_dns_record" "example" {
  account_name = "my_account"
  name         = "myalias"
  deployment   = "ss123456"
  ttl          = "300"
}
