package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhooksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_webhooks" "test" { account_name = "test_account_name" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_webhooks.test", "webhooks.#", "1"),
					resource.TestCheckResourceAttr("data.searchstax_webhooks.test", "webhooks.0.name", "webhookInput"),
				),
			},
		},
	})
}
