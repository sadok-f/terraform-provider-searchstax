package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRestoreStatusDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_restore_status" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  backup_id      = "27004"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_restore_status.test", "message", "Backup Restore in Progress"),
					resource.TestCheckResourceAttr("data.searchstax_restore_status.test", "status", "In Progress"),
				),
			},
		},
	})
}
