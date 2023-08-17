package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrivateVPCDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "searchstax_private_vpc" "test"{
							account_name="test_account_name"
						}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of deployments returned
					resource.TestCheckResourceAttr("data.searchstax_private_vpc.test", "private_vpc_list.#", "2"),
					// Verify the first deployment to ensure all attributes are set
					resource.TestCheckResourceAttr("data.searchstax_private_vpc.test", "private_vpc_list.0.name", "test-vpc"),
					resource.TestCheckResourceAttr("data.searchstax_private_vpc.test", "private_vpc_list.0.account", "test_account_name"),
					resource.TestCheckResourceAttr("data.searchstax_private_vpc.test", "private_vpc_list.0.region", "us-east1"),
					resource.TestCheckResourceAttr("data.searchstax_private_vpc.test", "private_vpc_list.0.status", "Active"),
					resource.TestCheckResourceAttr("data.searchstax_private_vpc.test", "private_vpc_list.0.address_space", "10.63.0.0/16"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.searchstax_private_vpc.test", "id", "placeholder"),
				),
			},
		},
	})
}
