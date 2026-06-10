resource "searchstax_deployment_rolling_restart" "example" {
  account_name   = "my_account"
  deployment_uid = "ss123456"
  solr           = true
  zookeeper      = false
}
