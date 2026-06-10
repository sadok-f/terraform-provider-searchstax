package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSRecordDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_dns_record" "test" {
  account_name = "test_account_name"
  name         = "myalias"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_dns_record.test", "name", "myalias"),
					resource.TestCheckResourceAttr("data.searchstax_dns_record.test", "deployment", "ss123456"),
					resource.TestCheckResourceAttr("data.searchstax_dns_record.test", "ttl", "300"),
				),
			},
		},
	})
}
