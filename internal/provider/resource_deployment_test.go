package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "searchstax_deployment" "test" {
  account_name             = "test_account_name"
  name                     = "SolrFromAPI"
  application              = "Solr"
  application_version      = "8.11.2"
  termination_lock         = "false"
  plan_type                = "DedicatedDeployment"
  plan                     = "NDC4-GCP-G"
  region_id                = "us-west-1"
  cloud_provider_id        = "gcp"
  num_additional_app_nodes = "0"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify returned values
					resource.TestCheckResourceAttr("searchstax_deployment.test", "name", "SolrFromAPI"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "application_version", "8.11.2"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "plan_type", "DedicatedDeployment"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "plan", "NDC4-GCP-G"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "region_id", "us-west-1"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "cloud_provider_id", "gcp"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "num_additional_app_nodes", "0"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("searchstax_deployment.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "searchstax_deployment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "test_account_name/ss123456",
				// The private_vpc attribute does not exist in the SearchStax
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"private_vpc"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
			resource "searchstax_deployment" "test" {
			 account_name             = "test_account_name"
			 name                     = "SolrFromAPI"
			 application              = "Solr"
			 application_version      = "8.11.2"
			 termination_lock         = "false"
			 plan_type                = "DedicatedDeployment"
			 plan                     = "NDC4-GCP-G"
			 region_id                = "us-west-1"
			 cloud_provider_id        = "gcp"
			 num_additional_app_nodes = "0"
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify deployment updated
					resource.TestCheckResourceAttr("searchstax_deployment.test", "uid", "ss123456"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "name", "SolrFromAPI"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "application_version", "8.11.2"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "plan_type", "DedicatedDeployment"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "plan", "NDC4-GCP-G"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "region_id", "us-west-1"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "cloud_provider_id", "gcp"),
					resource.TestCheckResourceAttr("searchstax_deployment.test", "num_additional_app_nodes", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
