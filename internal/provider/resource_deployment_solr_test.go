package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentSolrResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_deployment_solr" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  node           = "ss123456-1"
  action         = "start"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_deployment_solr.test", "action", "start"),
					resource.TestCheckResourceAttr("searchstax_deployment_solr.test", "node", "ss123456-1"),
				),
			},
		},
	})
}
