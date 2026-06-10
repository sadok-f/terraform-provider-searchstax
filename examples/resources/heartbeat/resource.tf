resource "searchstax_heartbeat" "example" {
  account_name   = "my_account"
  deployment_uid = "ss123456"
  name           = "uptime-check"
  host           = "https://example.com/health"
  interval       = "5m"
  max_alerts     = "3"
}
