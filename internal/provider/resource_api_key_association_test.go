package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAPIKeyAssociationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_api_key" "test" {
  account_name = "test_account_name"
  scope        = ["deployment.dedicateddeployment"]
}

resource "searchstax_api_key_association" "test" {
  account_name   = searchstax_api_key.test.account_name
  api_key        = searchstax_api_key.test.api_key
  deployment_uid = "ss123456"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_api_key_association.test", "deployment_uid", "ss123456"),
				),
			},
		},
	})
}
