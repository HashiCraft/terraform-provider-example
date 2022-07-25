package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccBlockDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.example_block.test", "x", "24"),
					resource.TestCheckResourceAttr("data.example_block.test", "y", "63"),
					resource.TestCheckResourceAttr("data.example_block.test", "z", "65"),
				),
			},
		},
	})
}

const testAccBlockDataSourceConfig = `
provider "example" {
  url = "http://localhost:8080"
}

resource "example_block" "brick" {
  x = 24
  y = 63
  z = 65
  material = "minecraft:stone"
}

data "example_block" "test" {
  id = resource.example_block.brick.id
}
`
