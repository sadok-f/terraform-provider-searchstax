resource "searchstax_ip_filter" "example" {
  account_name   = "my_account"
  deployment_uid = "ss123456"
  cidr_ip        = "203.0.113.0/24"
  description    = "Office network"
  services       = ["solr"]
}
