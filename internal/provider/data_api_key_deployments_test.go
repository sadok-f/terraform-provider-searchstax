package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAPIKeyDeploymentsDataSource(t *testing.T) {
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
}

data "searchstax_api_key_deployments" "test" {
  account_name = searchstax_api_key.test.account_name
  api_key      = searchstax_api_key.test.api_key
  depends_on   = [searchstax_api_key_association.test]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_api_key_deployments.test", "deployments.#", "1"),
					resource.TestCheckResourceAttr("data.searchstax_api_key_deployments.test", "deployments.0", "ss123456"),
				),
			},
		},
	})
}
