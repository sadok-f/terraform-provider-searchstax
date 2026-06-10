package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentBackupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_deployment_backup" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_deployment_backup.test", "backup_id", "27004"),
					resource.TestCheckResourceAttr("searchstax_deployment_backup.test", "id", "test_account_name/ss123456/27004"),
				),
			},
		},
	})
}
