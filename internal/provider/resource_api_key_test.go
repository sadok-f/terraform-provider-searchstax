package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAPIKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_api_key" "test" {
  account_name = "test_account_name"
  scope        = ["deployment.dedicateddeployment"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("searchstax_api_key.test", "api_key"),
					resource.TestCheckResourceAttr("searchstax_api_key.test", "scope.#", "1"),
				),
			},
		},
	})
}
