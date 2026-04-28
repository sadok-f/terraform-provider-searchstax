package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhookResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "searchstax_webhook" "test" {
  account_name = "test_account_name"
  webhook_id   = 7
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("searchstax_webhook.test", "webhook_id", "7"),
					resource.TestCheckResourceAttr("searchstax_webhook.test", "name", "webhookInput"),
				),
			},
		},
	})
}
