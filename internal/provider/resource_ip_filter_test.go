package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPFilterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_ip_filter" "test" {
  account_name   = "test_account_name"
  deployment_uid = "ss123456"
  cidr_ip        = "100.100.100.100/32"
  description    = "Added by API"
  services       = ["solr", "zk"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_ip_filter.test", "cidr_ip", "100.100.100.100/32"),
					resource.TestCheckResourceAttr("searchstax_ip_filter.test", "services.#", "2"),
				),
			},
		},
	})
}
