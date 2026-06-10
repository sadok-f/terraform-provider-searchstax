package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZookeeperConfigDownloadDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_zookeeper_config_download" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  name           = "film"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_zookeeper_config_download.test", "download", "mock"),
				),
			},
		},
	})
}
