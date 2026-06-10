resource "searchstax_backup_schedule" "example" {
  account_name   = "my_account"
  deployment_uid = "ss123456"
  days           = ["Monday", "Wednesday"]
  retention      = 7
  region_id      = "us-west-1"
  time           = "02:00"
}
