terraform {
  required_providers {
    searchstax = {
      source = "hashicorp.com/sadok-f/searchstax"
    }
  }
}

provider "searchstax" {
  username = var.ssx_username
  password = var.ssx_pwd
  host     = var.ssx_host
}

# --- Read-only: account & deployment context ---

data "searchstax_deployments" "all" {
  account_name = var.account_name
}

data "searchstax_deployment" "selected" {
  account_name   = var.account_name
  deployment_uid = var.deployment_uid
}

data "searchstax_deployment_health" "selected" {
  account_name   = var.account_name
  deployment_uid = var.deployment_uid
}

data "searchstax_plans" "available" {
  account_name = var.account_name
  application  = "Solr"
  plan_type    = "DedicatedPlan"
}

# --- Deployment management (optional: comment out if managing an existing deployment) ---

resource "searchstax_basic_auth" "example" {
  account_name   = var.account_name
  deployment_uid = var.deployment_uid
  enabled        = true
}

resource "searchstax_ip_filter" "office" {
  account_name   = var.account_name
  deployment_uid = var.deployment_uid
  cidr_ip        = var.office_cidr
  description    = "Office egress"
  services       = ["solr"]
}

resource "searchstax_tags_set" "example" {
  account_name   = var.account_name
  deployment_uid = var.deployment_uid
  tags           = var.deployment_tags
}

# --- Outputs ---

output "deployment_endpoint" {
  value = data.searchstax_deployment.selected.http_endpoint
}

output "deployment_health_status" {
  value = data.searchstax_deployment_health.selected.status
}

output "plan_names" {
  value = [for p in data.searchstax_plans.available.plans : p.name]
}
