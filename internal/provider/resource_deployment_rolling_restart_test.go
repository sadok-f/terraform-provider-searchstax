package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentRollingRestartResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_deployment_rolling_restart" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  solr           = true
  zookeeper      = false
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("searchstax_deployment_rolling_restart.test", "message"),
				),
			},
		},
	})
}
