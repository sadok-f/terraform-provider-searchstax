package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZookeeperConfigResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_zookeeper_config" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  name           = "test_config"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_zookeeper_config.test", "name", "test_config"),
					resource.TestCheckResourceAttr("searchstax_zookeeper_config.test", "id", "test_account_name/ss123456/test_config"),
				),
			},
		},
	})
}
