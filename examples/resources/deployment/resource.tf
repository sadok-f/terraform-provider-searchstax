resource "searchstax_deployment" "example" {
  account_name             = "my_account"
  name                     = "SolrFromTerraform"
  application              = "Solr"
  application_version      = "8.11.2"
  termination_lock         = false
  plan_type                = "DedicatedDeployment"
  plan                     = "NDC4-GCP-G"
  region_id                = "us-west-1"
  cloud_provider_id        = "gcp"
  num_additional_app_nodes = 0
}
