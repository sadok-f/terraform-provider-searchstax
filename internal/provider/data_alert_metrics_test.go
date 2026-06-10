package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlertMetricsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_alert_metrics" "test" {
  account_name = "test_account_name"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_alert_metrics.test", "metrics.#", "1"),
					resource.TestCheckResourceAttr("data.searchstax_alert_metrics.test", "metrics.0.metric", "system_cpu_usage"),
					resource.TestCheckResourceAttr("data.searchstax_alert_metrics.test", "metrics.0.unit", "percentage"),
				),
			},
		},
	})
}
