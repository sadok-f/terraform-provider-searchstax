resource "searchstax_deployment_solr" "example" {
  account_name   = "my_account"
  deployment_uid = "ss123456"
  node           = "ss123456-1"
  action         = "start"
}
