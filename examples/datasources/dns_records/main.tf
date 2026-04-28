data "searchstax_dns_records" "example" {
  account_name = var.account_name
}

output "dns_records" {
  value = data.searchstax_dns_records.example.records
}

