data "searchstax_auth_token" "example" {}

output "token" {
  value     = data.searchstax_auth_token.example.token
  sensitive = true
}

