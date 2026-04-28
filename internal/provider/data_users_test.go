package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUsersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_users" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_users.test", "users.#", "1"),
					resource.TestCheckResourceAttr("data.searchstax_users.test", "users.0.email", "user@company.com"),
					resource.TestCheckResourceAttr("data.searchstax_users.test", "users.0.role", "Admin"),
				),
			},
		},
	})
}
