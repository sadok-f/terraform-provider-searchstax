package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBackupScheduleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_backup_schedule" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  days           = ["mon", "wed", "fri"]
  time           = "07:00"
  retention      = 7
  region_id      = "us-west-1"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_backup_schedule.test", "retention", "7"),
					resource.TestCheckResourceAttr("searchstax_backup_schedule.test", "days.#", "3"),
					resource.TestCheckResourceAttrSet("searchstax_backup_schedule.test", "schedule_id"),
				),
			},
		},
	})
}
