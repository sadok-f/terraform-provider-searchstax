package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_deployment" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_deployment.test", "uid", "ss123456"),
					resource.TestCheckResourceAttr("data.searchstax_deployment.test", "name", "SolrFromAPI"),
					resource.TestCheckResourceAttr("data.searchstax_deployment.test", "application", "Solr"),
					resource.TestCheckResourceAttr("data.searchstax_deployment.test", "servers.#", "1"),
					resource.TestCheckResourceAttr("data.searchstax_deployment.test", "servers.0", "ss123456-1"),
					resource.TestCheckResourceAttr("data.searchstax_deployment.test", "subscription", "monthly"),
					resource.TestCheckResourceAttr("data.searchstax_deployment.test", "spec_jvm_heap_memory", "536870912"),
				),
			},
		},
	})
}
