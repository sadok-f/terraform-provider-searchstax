package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRestoreResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_restore" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  backup_id      = "27004"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_restore.test", "backup_id", "27004"),
					resource.TestCheckResourceAttrSet("searchstax_restore.test", "restore_id"),
					resource.TestCheckResourceAttrSet("searchstax_restore.test", "status"),
				),
			},
		},
	})
}
