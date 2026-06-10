package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlansDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "searchstax_plans" "test" {
  account_name = "test_account_name"
  application  = "Solr"
  plan_type    = "DedicatedPlan"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.searchstax_plans.test", "total_count", "2"),
					resource.TestCheckResourceAttr("data.searchstax_plans.test", "plans.#", "2"),
					resource.TestCheckResourceAttr("data.searchstax_plans.test", "plans.0.application", "Solr"),
				),
			},
		},
	})
}
