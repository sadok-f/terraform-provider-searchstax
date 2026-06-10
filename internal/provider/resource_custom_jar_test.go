package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomJarResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_custom_jar" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  name           = "test.jar"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_custom_jar.test", "name", "test.jar"),
					resource.TestCheckResourceAttr("searchstax_custom_jar.test", "id", "test_account_name/ss123456/test.jar"),
				),
			},
		},
	})
}
