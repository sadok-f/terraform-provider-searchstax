package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_tags" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_tags.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("data.searchstax_tags.test", "tags.0", "demo"),
					resource.TestCheckResourceAttr("data.searchstax_tags.test", "tags.1", "test"),
				),
			},
		},
	})
}
