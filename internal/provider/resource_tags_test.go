package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagsSetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_tags_set" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  tags           = ["demo", "test"]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_tags_set.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("searchstax_tags_set.test", "tags.0", "demo"),
					resource.TestCheckResourceAttr("searchstax_tags_set.test", "tags.1", "test"),
					resource.TestCheckResourceAttr("searchstax_tags_set.test", "id", "test_account_name/ss123456"),
				),
			},
		},
	})
}
