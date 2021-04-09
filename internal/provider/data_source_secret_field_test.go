package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSecretField(t *testing.T) {
	t.Skip("data source not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretField,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.tss_secret_field.foo", "number", regexp.MustCompile("^1")),
				),
			},
		},
	})
}

const testAccDataSourceSecretField = `
data "tss_secret_field" "foo" {
  number = 1
  slug = "bar"
}
`
