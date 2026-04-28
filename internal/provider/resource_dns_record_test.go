package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_dns_record" "test" {
  account_name = "test_account_name"
  name         = "myalias"
  deployment   = "ss123456"
  ttl          = "300"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_dns_record.test", "name", "myalias"),
					resource.TestCheckResourceAttr("searchstax_dns_record.test", "deployment", "ss123456"),
				),
			},
		},
	})
}
