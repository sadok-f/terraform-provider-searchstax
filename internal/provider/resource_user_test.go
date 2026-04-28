package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_user" "test" {
  email      = "user@company.com"
  role       = "Admin"
  first_name = "Mock"
  last_name  = "User"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_user.test", "email", "user@company.com"),
					resource.TestCheckResourceAttr("searchstax_user.test", "role", "Admin"),
				),
			},
			{
				ResourceName:      "searchstax_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "user@company.com",
				ImportStateVerifyIgnore: []string{
					"new_password",
				},
			},
		},
	})
}
