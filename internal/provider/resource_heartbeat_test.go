package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHeartbeatResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_heartbeat" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  name           = "mock-heartbeat"
  host           = "*"
  interval       = "5"
  max_alerts     = "5"
  email          = ["user@company.com"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_heartbeat.test", "name", "mock-heartbeat"),
					resource.TestCheckResourceAttr("searchstax_heartbeat.test", "host", "*"),
					resource.TestCheckResourceAttr("searchstax_heartbeat.test", "interval", "5"),
					resource.TestCheckResourceAttrSet("searchstax_heartbeat.test", "heartbeat_id"),
				),
			},
		},
	})
}
