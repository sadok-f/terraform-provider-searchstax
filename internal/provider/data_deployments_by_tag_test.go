package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentsByTagDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_deployments_by_tag" "test" {
  account_name = "test_account_name"
  tags           = ["demo"]
  operator       = "OR"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_deployments_by_tag.test", "results.#", "1"),
					resource.TestCheckResourceAttr("data.searchstax_deployments_by_tag.test", "results.0.deployment", "ss123456"),
					resource.TestCheckResourceAttr("data.searchstax_deployments_by_tag.test", "results.0.tags.#", "2"),
				),
			},
		},
	})
}
