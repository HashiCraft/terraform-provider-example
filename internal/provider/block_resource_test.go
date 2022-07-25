package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/client"
	"github.com/sethvargo/go-retry"
)

func TestAccBlockResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBlockResourceConfig("minecraft:stone"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("example_block.brick", "x", "24"),
					resource.TestCheckResourceAttr("example_block.brick", "y", "63"),
					resource.TestCheckResourceAttr("example_block.brick", "z", "64"),
					resource.TestCheckResourceAttr("example_block.brick", "id", "24_63_64"),
					customCheckStep,
				),
			},
		},
	})
}

func customCheckStep(s *terraform.State) error {
	c := client.NewClient("http://localhost:8080")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return retry.Constant(ctx, 1*time.Second, func(ctx context.Context) error {
		b, err := c.GetBlock(24, 63, 64)

		if err != nil {
			return retry.RetryableError(err)
		}

		if b == nil {
			return retry.RetryableError(fmt.Errorf("Unable to find block"))
		}

		return nil
	})
}

func testAccBlockResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
provider "example" {
  url = "http://localhost:8080"
}

resource "example_block" "brick" {
  x = 24
  y = 63
  z = 64
  material = "%s"
}
`, configurableAttribute)
}
