package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "searchstax_deployments" "test"{
							account_name="test_account_name"
						}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of deployments returned
					resource.TestCheckResourceAttr("data.searchstax_deployments.test", "deployments_list.#", "2"),
					// Verify the first deployment to ensure all attributes are set
					resource.TestCheckResourceAttr("data.searchstax_deployments.test", "deployments_list.0.name", "ListByAPI"),
					resource.TestCheckResourceAttr("data.searchstax_deployments.test", "deployments_list.0.uid", "ss123456"),
					resource.TestCheckResourceAttr("data.searchstax_deployments.test", "deployments_list.0.application", "Solr"),
					resource.TestCheckResourceAttr("data.searchstax_deployments.test", "deployments_list.0.application_version", "8.11.2"),
					resource.TestCheckResourceAttr("data.searchstax_deployments.test", "deployments_list.0.tier", "Gold"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.searchstax_deployments.test", "id", "placeholder"),
				),
			},
		},
	})
}
