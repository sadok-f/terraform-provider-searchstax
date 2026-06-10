package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_deployment_user" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  username       = "demoSolr"
  password       = "test123"
  role           = "Admin"
}`,
				ResourceName:  "searchstax_deployment_user.test",
				ImportState:   true,
				ImportStateId: "test_account_name/ss123456/demoSolr",
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_deployment_user.test", "username", "demoSolr"),
					resource.TestCheckResourceAttr("searchstax_deployment_user.test", "role", "Admin"),
					resource.TestCheckResourceAttr("searchstax_deployment_user.test", "id", "test_account_name/ss123456/demoSolr"),
				),
			},
			{
				Config: providerConfig + `
resource "searchstax_deployment_user" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  username       = "demoSolr"
  password       = "updated-password"
  role           = "Admin"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_deployment_user.test", "password", "updated-password"),
				),
			},
		},
	})
}
